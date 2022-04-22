package znet

import (
	"net"
)

// ClientConnection 客户端使用的连接对象
type ClientConnection struct {
	*BaseConnection
}

// NewClientConnection 新建客户端使用的连接对象
func NewClientConnection(c *Client, conn *net.TCPConn) *ClientConnection {
	bc := NewBaseConnection(c, conn)
	return &ClientConnection{
		BaseConnection: bc,
	}
}
