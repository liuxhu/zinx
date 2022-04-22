package znet

import "github.com/liuxhu/zinx/ziface"

// ServerOption 服务端可选参数
type ServerOption func(s *Server)

// WithConnManager 传入自定义实现的conn管理对象
func WithConnManager(connMgr ziface.IConnManager) ServerOption {
	return func(s *Server) {
		s.ConnMgr = connMgr
	}
}

// WithMaxConn 传入允许的最大连接数
func WithMaxConn(maxConn int) ServerOption {
	return func(s *Server) {
		s.maxConn = maxConn
	}
}
