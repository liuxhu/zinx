package znet

import (
	"time"

	"github.com/liuxhu/zinx/ziface"
)

// BaseService 基础服务对象
type BaseService struct {
	ipVersion     string
	addr          string
	onConnStart   func(ziface.IConnection)
	onConnStop    func(ziface.IConnection)
	msgHandler    ziface.IMsgHandle
	packet        ziface.Packet
	readDeadline  time.Duration
	writeDeadline time.Duration
}

// NewBaseService 新建基础服务对象
func NewBaseService(opts ...Option) ziface.IService {
	s := &BaseService{
		ipVersion:     defaultIPVersion,
		addr:          defaultAddr,
		msgHandler:    NewMsgHandle(),
		packet:        NewDataPack(),
		readDeadline:  defaultReadDeadline,
		writeDeadline: defaultWriteDeadline,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start ...
func (s *BaseService) Start() {
	panic("implement me")
}

// Stop ...
func (s *BaseService) Stop() {}

// Serve ...
func (s *BaseService) Serve() {}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *BaseService) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.onConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *BaseService) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.onConnStop = hookFunc
}

// CallOnConnStart 调用连接OnConnStart Hook函数
func (s *BaseService) CallOnConnStart(conn ziface.IConnection) {
	if s.onConnStart != nil {
		s.onConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook函数
func (s *BaseService) CallOnConnStop(conn ziface.IConnection) {
	if s.onConnStop != nil {
		s.onConnStop(conn)
	}
}

// IPVersion ...
func (s *BaseService) IPVersion() string {
	return s.ipVersion
}

// Addr ...
func (s *BaseService) Addr() string {
	return s.addr
}

// AddRouter 添加路由
func (s *BaseService) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

// Packet 获取packet对象
func (s *BaseService) Packet() ziface.Packet {
	return s.packet
}

// MsgHandler ...
func (s *BaseService) MsgHandler() ziface.IMsgHandle {
	return s.msgHandler
}

// ReadDeadline ...
func (s *BaseService) ReadDeadline() time.Duration {
	return s.readDeadline
}

// WriteDeadline ...
func (s *BaseService) WriteDeadline() time.Duration {
	return s.writeDeadline
}
