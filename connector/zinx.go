package connector

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/tcpface"
	"github.com/jqiris/kungfu/tcpserver/zinx"
	"github.com/jqiris/kungfu/treaty"
)

type ZinxConnector struct {
	ServerId              string
	Server                *treaty.Server
	RpcX                  rpcx.RpcConnector
	EventHandlerSelf      func(req []byte) []byte //处理自己的事件
	EventHandlerBroadcast func(req []byte) []byte //处理广播事件
	ClientServer          tcpface.IServer         //zinx server
	RouteHandler          func(s tcpface.IServer) //注册路由
}

func (b *ZinxConnector) Init() {
	//find the  server config
	if serverConf := helper.FindServerConfig(config.GetServersConf(), b.GetServerId()); serverConf == nil {
		logger.Fatal("ZinxConnector can't find the server config")
	} else {
		b.Server = serverConf
	}
	//init the rpcx
	b.RpcX = rpcx.NewRpcConnector(config.GetRpcXConf())
	//run the front server
	b.ClientServer = zinx.NewServer(b.Server)
	b.RouteHandler(b.ClientServer)
	go b.ClientServer.Serve()

	logger.Infoln("init the connector:", b.ServerId)
}

func (b *ZinxConnector) AfterInit() {
	//Subscribe event
	if err := b.RpcX.Subscribe(b.Server, func(req []byte) []byte {
		logger.Infof("ZinxConnector Subscribe received: %+v", req)
		return b.EventHandlerSelf(req)
	}); err != nil {
		logger.Error(err)
	}
	if err := b.RpcX.SubscribeConnector(func(req []byte) []byte {
		logger.Infof("ZinxConnector SubscribeConnector received: %+v", req)
		return b.EventHandlerBroadcast(req)
	}); err != nil {
		logger.Error(err)
	}
	//register the service
	if err := discover.Register(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *ZinxConnector) BeforeShutdown() {
	//unregister the service
	if err := discover.UnRegister(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *ZinxConnector) Shutdown() {
	//stop the server
	if b.ClientServer != nil {
		b.ClientServer.Stop()
	}
	logger.Infoln("stop the connector:", b.ServerId)
}

func (b *ZinxConnector) GetServer() *treaty.Server {
	return b.Server
}

func (b *ZinxConnector) RegEventHandlerSelf(handler func(req []byte) []byte) { //注册自己事件处理器
	b.EventHandlerSelf = handler
}

func (b *ZinxConnector) RegEventHandlerBroadcast(handler func(req []byte) []byte) { //注册广播事件处理器
	b.EventHandlerBroadcast = handler
}
func (b *ZinxConnector) SetServerId(serverId string) {
	b.ServerId = serverId
}

func (b *ZinxConnector) GetServerId() string {
	return b.ServerId
}

func (b *ZinxConnector) SetRouteHandler(handler func(s tcpface.IServer)) {
	b.RouteHandler = handler
}
