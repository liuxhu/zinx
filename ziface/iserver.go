// Package ziface 主要提供zinx全部抽象层接口定义.
// 包括:
//		IServer 服务mod接口
//		IRouter 路由mod接口
//		IConnection 连接mod层接口
//      IMessage 消息mod接口
//		IDataPack 消息拆解接口
//      IMsgHandler 消息处理及协程池接口
//
// 当前文件描述:
// @Title  iserver.go
// @Description  提供Server抽象层全部接口声明
// @Author  Aceld - Thu Mar 11 10:32:29 CST 2019
package ziface

// IServer 服务端接口
type IServer interface {
	IService

	GetConnMgr() IConnManager //得到链接管理
}
