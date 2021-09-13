package connector

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/tcpface"
	"github.com/jqiris/kungfu/tcpserver"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/kungfu/utils"
)

type TcpConnector struct {
	ServerId              string
	Server                *treaty.Server
	RpcX                  rpcx.RpcServer
	EventJsonSelf         rpcx.CallbackFunc       //处理自己的json事件
	EventHandlerSelf      rpcx.CallbackFunc       //处理自己的事件
	EventHandlerBroadcast rpcx.CallbackFunc       //处理广播事件
	ClientServer          tcpface.IServer         //zinx server
	RouteHandler          func(s tcpface.IServer) //注册路由
}

func (b *TcpConnector) Init() {
	//find the  server config
	if serverConf := utils.FindServerConfig(config.GetServersConf(), b.GetServerId()); serverConf == nil {
		logger.Fatal("NanoConnector can't find the server config")
	} else {
		b.Server = serverConf
	}
	//init the rpcx
	b.RpcX = rpcx.NewRpcServer(config.GetRpcXConf(), b.Server)
	//run the front server
	b.ClientServer = tcpserver.NewServer(b.Server)
	b.RouteHandler(b.ClientServer)
	go b.ClientServer.Serve()

	logger.Info("init the connector:", b.ServerId)
}

func (b *TcpConnector) AfterInit() {
	//Subscribe event
	if err := b.RpcX.Subscribe(b.Server, func(req *rpcx.RpcMsg) []byte {
		logger.Infof("NanoConnector Subscribe received: %+v", req)
		return b.EventHandlerSelf(req)
	}); err != nil {
		logger.Error(err)
	}
	//Subscribe event
	if err := b.RpcX.SubscribeJson(b.Server, func(req *rpcx.RpcMsg) []byte {
		//logger.Infof("BaseBackEnd Subscribe received: %+v", req)
		return b.EventJsonSelf(req)
	}); err != nil {
		logger.Error(err)
	}
	if err := b.RpcX.SubscribeConnector(func(req *rpcx.RpcMsg) []byte {
		logger.Infof("NanoConnector SubscribeConnector received: %+v", req)
		return b.EventHandlerBroadcast(req)
	}); err != nil {
		logger.Error(err)
	}
	//register the service
	if err := discover.Register(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *TcpConnector) BeforeShutdown() {
	//unregister the service
	if err := discover.UnRegister(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *TcpConnector) Shutdown() {
	//stop the server
	if b.ClientServer != nil {
		b.ClientServer.Stop()
	}
	logger.Info("stop the connector:", b.ServerId)
}

func (b *TcpConnector) GetServer() *treaty.Server {
	return b.Server
}

func (b *TcpConnector) RegEventHandlerSelf(handler rpcx.CallbackFunc) { //注册自己事件处理器
	b.EventHandlerSelf = handler
}

func (b *TcpConnector) RegEventHandlerBroadcast(handler rpcx.CallbackFunc) { //注册广播事件处理器
	b.EventHandlerBroadcast = handler
}
func (b *TcpConnector) SetServerId(serverId string) {
	b.ServerId = serverId
}

func (b *TcpConnector) GetServerId() string {
	return b.ServerId
}

func (b *TcpConnector) SetRouteHandler(handler func(s tcpface.IServer)) {
	b.RouteHandler = handler
}
