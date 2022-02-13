package connector

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
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
	if b.Server == nil {
		panic("服务配置信息不能为空")
		return
	}
	//赋值id
	b.ServerId = b.Server.ServerId
	b.ConnectorConf = config.GetConnectorConf()
	//init the rpcx
	b.RpcX = rpcx.NewRpcServer(config.GetRpcXConf(), b.Server)
}

func (b *HttpConnector) AfterInit() {
	if b.Server == nil {
		panic("服务配置信息不能为空")
		return
	}
	if b.EventJsonSelf == nil {
		panic("EventJsonSelf不能为空")
		return
	}
	if b.EventHandlerSelf == nil {
		panic("EventHandlerSelf不能为空")
		return
	}
	if b.EventHandlerBroadcast == nil {
		panic("EventHandlerBroadcast不能为空")
		return
	}
	builder := rpcx.NewRpcSubscriber(b.Server).SetCodeType(rpcx.CodeTypeProto).SetCallback(func(req *rpcx.RpcMsg) []byte {
		return b.EventHandlerSelf(req)
	})
	//Subscribe event
	if err := b.RpcX.Subscribe(*builder); err != nil {
		logger.Error(err)
	}
	builder = builder.SetSuffix("json").SetCodeType(rpcx.CodeTypeJson).SetCallback(func(req *rpcx.RpcMsg) []byte {
		return b.EventJsonSelf(req)
	})
	//Subscribe event
	if err := b.RpcX.Subscribe(*builder); err != nil {
		logger.Error(err)
	}
	builder = builder.SetSuffix(rpcx.DefaultSuffix).SetCodeType(rpcx.CodeTypeProto).SetCallback(func(req *rpcx.RpcMsg) []byte {
		return b.EventHandlerBroadcast(req)
	})
	if err := b.RpcX.SubscribeConnector(*builder); err != nil {
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
