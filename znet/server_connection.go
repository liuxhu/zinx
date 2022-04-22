package znet

import (
	"net"

	"github.com/liuxhu/zinx/ziface"
)

// ServerConnection 服务端使用的连接对象
type ServerConnection struct {
	*BaseConnection

	//当前Conn属于哪个Server
	TCPServer ziface.IServer
	ConnID    uint32
}

// NewServerConnection 创建服务端使用的连接对象
func NewServerConnection(s ziface.IServer, conn *net.TCPConn, connID uint32) *ServerConnection {
	// 初始化Conn属性
	bc := NewBaseConnection(s, conn)
	sc := &ServerConnection{
		BaseConnection: bc,
		TCPServer:      s,
		ConnID:         connID,
	}

	// 将新创建的Conn添加到链接管理中
	sc.TCPServer.GetConnMgr().Add(sc)

	return sc
}

//Start 启动连接，让当前连接开始工作
func (c *ServerConnection) Start() {
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TCPServer.CallOnConnStart(c)
	c.BaseConnection.Start()
}

func (c *ServerConnection) finalizer() {
	c.TCPServer.GetConnMgr().Remove(c)

	c.BaseConnection.finalizer()
}
