package znet

import "time"

// ClientOption 客户端可选参数
type ClientOption func(*Client)

// WithReconnectWaitTime 传入重连等待时间
func WithReconnectWaitTime(t time.Duration) ClientOption {
	return func(c *Client) {
		c.reconnectWaitTime = t
	}
}

// WithMaxReconnectCount 传入最大重连次数
func WithMaxReconnectCount(count uint32) ClientOption {
	return func(c *Client) {
		c.maxReconnectCount = count
	}
}
