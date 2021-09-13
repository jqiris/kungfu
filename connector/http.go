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
	EventJsonSelf         rpcx.CallbackFunc //处理自己的json事件
	EventHandlerSelf      rpcx.CallbackFunc //处理自己的事件
	EventHandlerBroadcast rpcx.CallbackFunc //处理广播事件
	ConnectorConf         config.ConnectorConf
}

func (b *HttpConnector) Init() {
	//find the  server config
	if serverConf := utils.FindServerConfig(config.GetServersConf(), b.GetServerId()); serverConf == nil {
		logger.Fatal("HttpConnector can't find the server config")
	} else {
		b.Server = serverConf
		b.ConnectorConf = config.GetConnectorConf()
	}
	//init the rpcx
	b.RpcX = rpcx.NewRpcServer(config.GetRpcXConf(), b.Server)
}

func (b *HttpConnector) AfterInit() {
	//Subscribe event
	if err := b.RpcX.Subscribe(b.Server, func(req *rpcx.RpcMsg) []byte {
		logger.Infof("HttpConnector Subscribe received: %+v", req)
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
		logger.Infof("HttpConnector SubscribeConnector received: %+v", req)
		return b.EventHandlerBroadcast(req)
	}); err != nil {
		logger.Error(err)
	}
	//register the service
	if err := discover.Register(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *HttpConnector) BeforeShutdown() {
	//unregister the service
	if err := discover.UnRegister(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *HttpConnector) Shutdown() {
	logger.Info("stop the connector:", b.ServerId)
}

func (b *HttpConnector) GetServer() *treaty.Server {
	return b.Server
}
func (b *HttpConnector) RegEventJsonSelf(handler rpcx.CallbackFunc) { //注册自己事件处理器
	b.EventJsonSelf = handler
}
func (b *HttpConnector) RegEventHandlerSelf(handler rpcx.CallbackFunc) { //注册自己事件处理器
	b.EventHandlerSelf = handler
}

func (b *HttpConnector) RegEventHandlerBroadcast(handler rpcx.CallbackFunc) { //注册广播事件处理器
	b.EventHandlerBroadcast = handler
}
func (b *HttpConnector) SetServerId(serverId string) {
	b.ServerId = serverId
}

func (b *HttpConnector) GetServerId() string {
	return b.ServerId
}
