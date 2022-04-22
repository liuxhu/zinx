package znet

import (
	"time"

	"github.com/liuxhu/zinx/ziface"
)

// Option 可选参数
type Option func(s *BaseService)

// WithPacket 只要实现Packet 接口可自由实现数据包解析格式，如果没有则使用默认解析格式
func WithPacket(pack ziface.Packet) Option {
	return func(s *BaseService) {
		s.packet = pack
	}
}

// WithAddr 传入监听地址
func WithAddr(addr string) Option {
	return func(s *BaseService) {
		s.addr = addr
	}
}

// WithReadDeadline 传入读取超时时间
func WithReadDeadline(t time.Duration) Option {
	return func(s *BaseService) {
		s.readDeadline = t
	}
}

// WithWriteDeadline 传入写入超时时间
func WithWriteDeadline(t time.Duration) Option {
	return func(s *BaseService) {
		s.writeDeadline = t
	}
}
