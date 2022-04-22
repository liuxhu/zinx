package znet

import (
	"fmt"
	"net"
	"time"

	"github.com/liuxhu/zinx/ziface"
)

// Client 客户端对象
type Client struct {
	ziface.IService

	reconnectWaitTime time.Duration
	maxReconnectCount uint32
	reconnectCount    uint32
}

// NewClient 新建一个客户端对象
func NewClient(cOpts []ClientOption, opts ...Option) ziface.IClient {
	c := &Client{
		IService:       NewBaseService(opts...),
		reconnectCount: 0,
	}

	for _, opt := range cOpts {
		opt(c)
	}

	return c
}

// Start 启动
func (c *Client) Start() {
	addr, err := net.ResolveTCPAddr(c.IPVersion(), c.Addr())
	if err != nil {
		fmt.Println("resolve tcp addr err: ", err)
		return
	}

	conn, err := net.DialTCP(c.IPVersion(), nil, addr)
	if err != nil {
		fmt.Printf("dial:%s err: %v\n", addr.String(), err)
		return
	}

	dealConn := NewClientConnection(c, conn)
	go c.CallOnConnStart(dealConn)

	dealConn.Start()
}

// Serve ...
func (c *Client) Serve() {
	addr, err := net.ResolveTCPAddr(c.IPVersion(), c.Addr())
	if err != nil {
		fmt.Println("resolve tcp addr err: ", err)
		return
	}

	for {
		conn, err := net.DialTCP(c.IPVersion(), nil, addr)
		if err != nil {
			c.reconnectCount++
			if c.maxReconnectCount > 0 && c.reconnectCount >= c.maxReconnectCount {
				fmt.Println("reconnectCount over than maxReconnectCount, reconnect failed")
				break
			}

			time.Sleep(c.reconnectWaitTime)
			continue
		}

		c.reconnectCount = 0

		dealConn := NewClientConnection(c, conn)
		go c.CallOnConnStart(dealConn)

		dealConn.Start()
	}
}
