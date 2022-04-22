package znet

import (
	"fmt"
	"net"

	"github.com/liuxhu/zinx/ziface"
)

const (
	defaultIPVersion = "tcp4"
	defaultAddr      = "0.0.0.0:3000"
	defaultMaxConn   = 12000

	defaultWorkerPoolSize   = 10
	defaultMaxWorkerTaskLen = 1024
	defaultMaxMsgChanLen    = 1024
	defaultMaxPacketSize    = 0

	// 默认读写超时时间(0为不限制)
	defaultReadDeadline  = 0
	defaultWriteDeadline = 0
)

//Server 接口实现，定义一个Server服务类
type Server struct {
	ziface.IService

	ConnMgr ziface.IConnManager

	maxConn int
}

// NewServer 创建一个服务器句柄
func NewServer(sOpts []ServerOption, opts ...Option) ziface.IServer {
	s := &Server{
		IService: NewBaseService(opts...),
		ConnMgr:  NewConnManager(),
		maxConn:  defaultMaxConn,
	}

	for _, opt := range sOpts {
		opt(s)
	}

	return s
}

//============== 实现 ziface.IServer 里的全部接口方法 ========

//Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server ,listenner at Addr: %s is starting\n", s.Addr())

	//开启一个go去做服务端Linster业务
	go func() {
		//0 启动worker工作池机制
		s.MsgHandler().StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion(), s.Addr())
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion(), addr)
		if err != nil {
			panic(err)
		}

		//已经监听成功
		fmt.Println("start Zinx server success, now listenning...")

		//TODO server.go 应该有一个自动生成ID的方法
		var cID uint32
		cID = 0

		//3 启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())

			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= s.maxConn {
				conn.Close()
				continue
			}

			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewServerConnection(s, conn, cID)
			cID++

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}

//Stop 停止服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server")

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

//Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

//GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}
