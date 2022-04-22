package ziface

import "time"

// IService 定义服务接口
type IService interface {
	Start()
	Stop()
	Serve()
	AddRouter(uint32, IRouter)
	SetOnConnStart(func(IConnection)) //设置该Server的连接创建时Hook函数
	SetOnConnStop(func(IConnection))  //设置该Server的连接断开时的Hook函数
	CallOnConnStart(IConnection)      //调用连接OnConnStart Hook函数
	CallOnConnStop(IConnection)       //调用连接OnConnStop Hook函数
	IPVersion() string
	Addr() string
	Packet() Packet
	MsgHandler() IMsgHandle
	ReadDeadline() time.Duration
	WriteDeadline() time.Duration
}
