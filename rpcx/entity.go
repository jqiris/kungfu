package rpcx

import "github.com/jqiris/kungfu/treaty"

// ServerEntity server entity
type ServerEntity interface {
	Init()                                         //初始化
	AfterInit()                                    //初始化后执行操作
	BeforeShutdown()                               //服务关闭前操作
	Shutdown()                                     //服务关闭操作
	GetServer() *treaty.Server                     //获取服务
	RegEventHandlerSelf(handler CallbackFunc)      //注册自己事件处理器
	RegEventHandlerBroadcast(handler CallbackFunc) //注册广播事件处理器
	SetServerId(serverId string)                   //设置SetServerId
	GetServerId() string                           //获取serverId
}
