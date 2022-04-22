package znet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/liuxhu/zinx/ziface"
)

// BaseConnection 基础连接对象
type BaseConnection struct {
	s ziface.IService

	Conn   *net.TCPConn
	ConnID uint32

	// 取消相关
	ctx         context.Context
	cancel      context.CancelFunc
	msgBuffChan chan []byte

	// 链接属性
	sync.RWMutex
	property     map[string]interface{}
	propertyLock sync.Mutex
	isClosed     bool
}

// NewBaseConnection 创建基础连接对象
func NewBaseConnection(s ziface.IService, conn *net.TCPConn) *BaseConnection {
	//初始化Conn属性
	c := &BaseConnection{
		s:           s,
		Conn:        conn,
		isClosed:    false,
		msgBuffChan: make(chan []byte, defaultMaxMsgChanLen),
		property:    nil,
	}

	return c
}

// StartWriter 写消息Goroutine， 用户将数据发送给客户端
func (c *BaseConnection) StartWriter() {
	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if c.s.WriteDeadline() > 0 {
					_ = c.Conn.SetWriteDeadline(time.Now().Add(c.s.WriteDeadline()))
				}
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *BaseConnection) StartReader() {
	defer c.Stop()

	// 创建拆包解包的对象
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			//读取客户端的Msg head
			headData := make([]byte, c.s.Packet().GetHeadLen())
			if c.s.ReadDeadline() > 0 {
				_ = c.Conn.SetReadDeadline(time.Now().Add(c.s.ReadDeadline()))
			}
			if _, err := io.ReadFull(c.Conn, headData); err != nil {
				fmt.Println("read msg head error ", err)
				return
			}

			//拆包，得到msgID 和 datalen 放在msg中
			msg, err := c.s.Packet().Unpack(headData)
			if err != nil {
				fmt.Println("unpack error ", err)
				return
			}

			//根据 dataLen 读取 data，放在msg.Data中
			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if c.s.ReadDeadline() > 0 {
					_ = c.Conn.SetReadDeadline(time.Now().Add(c.s.ReadDeadline()))
				}
				if _, err := io.ReadFull(c.Conn, data); err != nil {
					fmt.Println("read msg data error ", err)
					return
				}
			}
			msg.SetData(data)

			//得到当前客户端请求的Request数据
			req := Request{
				conn: c,
				msg:  msg,
			}

			if defaultWorkerPoolSize > 0 {
				//已经启动工作池机制，将消息交给Worker处理
				c.s.MsgHandler().SendMsgToTaskQueue(&req)
			} else {
				//从绑定好的消息和对应的处理方法中执行对应的Handle方法
				go c.s.MsgHandler().DoMsgHandler(&req)
			}
		}
	}
}

// Start 启动连接，让当前连接开始工作
func (c *BaseConnection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

	select {
	case <-c.ctx.Done():
		c.finalizer()
		return
	}
}

// Stop 停止连接，结束当前连接状态M
func (c *BaseConnection) Stop() {
	c.cancel()
}

// GetTCPConnection 从当前连接获取原始的socket TCPConn
func (c *BaseConnection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接ID
func (c *BaseConnection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端地址信息
func (c *BaseConnection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *BaseConnection) SendMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}

	//将data封包，并且发送
	msg, err := c.s.Packet().Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg ID = ", msgID)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	_, err = c.Conn.Write(msg)
	return err
}

// SendBuffMsg  发生BuffMsg
func (c *BaseConnection) SendBuffMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()

	if c.isClosed {
		return errors.New("Connection closed when send buff msg")
	}

	//将data封包，并且发送
	msg, err := c.s.Packet().Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg ID = ", msgID)
		return errors.New("Pack error msg ")
	}

	// 发送超时
	select {
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- msg:
		return nil
	}
	//写回客户端
	//c.msgBuffChan <- msg

	return nil
}

// SetProperty 设置链接属性
func (c *BaseConnection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *BaseConnection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

// RemoveProperty 移除链接属性
func (c *BaseConnection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

// Context 返回ctx，用于用户自定义的go程获取连接退出状态
func (c *BaseConnection) Context() context.Context {
	return c.ctx
}

func (c *BaseConnection) finalizer() {
	c.s.CallOnConnStop(c)

	c.Lock()
	defer c.Unlock()

	//如果当前链接已经关闭
	if c.isClosed {
		return
	}

	// 关闭socket链接
	_ = c.Conn.Close()

	//关闭该链接全部管道
	close(c.msgBuffChan)
	//设置标志位
	c.isClosed = true
}
