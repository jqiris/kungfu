package connector

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/kungfu/utils"
)

type HttpConnector struct {
	ServerId              string
	Server                *treaty.Server
	RpcX                  rpcx.RpcServer
	EventHandlerSelf      rpcx.CallbackFunc //处理自己的事件
	EventHandlerBroadcast rpcx.CallbackFunc //处理广播事件
	ConnectorConf         config.ConnectorConf
}

func (g *HttpConnector) Init() {
	//find the  server config
	if serverConf := utils.FindServerConfig(config.GetServersConf(), g.GetServerId()); serverConf == nil {
		logger.Fatal("HttpConnector can't find the server config")
	} else {
		g.Server = serverConf
		g.ConnectorConf = config.GetConnectorConf()
	}
	//init the rpcx
	g.RpcX = rpcx.NewRpcServer(config.GetRpcXConf(), g.Server)
}

func (g *HttpConnector) AfterInit() {
	//Subscribe event
	if err := g.RpcX.Subscribe(g.Server, func(req *rpcx.RpcMsg) []byte {
		logger.Infof("HttpConnector Subscribe received: %+v", req)
		return g.EventHandlerSelf(req)
	}); err != nil {
		logger.Error(err)
	}
	if err := g.RpcX.SubscribeConnector(func(req *rpcx.RpcMsg) []byte {
		logger.Infof("HttpConnector SubscribeConnector received: %+v", req)
		return g.EventHandlerBroadcast(req)
	}); err != nil {
		logger.Error(err)
	}
	//register the service
	if err := discover.Register(g.Server); err != nil {
		logger.Error(err)
	}
}

func (g *HttpConnector) BeforeShutdown() {
	//unregister the service
	if err := discover.UnRegister(g.Server); err != nil {
		logger.Error(err)
	}
}

func (g *HttpConnector) Shutdown() {
	logger.Info("stop the connector:", g.ServerId)
}

func (g *HttpConnector) GetServer() *treaty.Server {
	return g.Server
}

func (g *HttpConnector) RegEventHandlerSelf(handler rpcx.CallbackFunc) { //注册自己事件处理器
	g.EventHandlerSelf = handler
}

func (g *HttpConnector) RegEventHandlerBroadcast(handler rpcx.CallbackFunc) { //注册广播事件处理器
	g.EventHandlerBroadcast = handler
}
func (g *HttpConnector) SetServerId(serverId string) {
	g.ServerId = serverId
}

func (g *HttpConnector) GetServerId() string {
	return g.ServerId
}
