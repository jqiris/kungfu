package backend

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
)

//BaseBackEnd 后盾服务
type BaseBackEnd struct {
	ServerId              string
	Server                *treaty.Server
	RpcX                  rpcx.RpcServer
	EventHandlerSelf      rpcx.CallbackFunc //处理自己的事件
	EventHandlerBroadcast rpcx.CallbackFunc //处理广播事件
}

func (b *BaseBackEnd) Init() {
	//find the  server config
	if b.Server = helper.FindServerConfig(config.GetServersConf(), b.GetServerId()); b.Server == nil {
		logger.Fatal("BaseBackEnd can find the server config")
	}
	//init the rpcx
	b.RpcX = rpcx.NewRpcServer(config.GetRpcXConf())
	logger.Infoln("init the backend:", b.ServerId)
}

func (b *BaseBackEnd) AfterInit() {
	//Subscribe event
	if err := b.RpcX.Subscribe(b.Server, func(coder *rpcx.RpcEncoder, req *rpcx.RpcMsg) []byte {
		//logger.Infof("BaseBackEnd Subscribe received: %+v", req)
		return b.EventHandlerSelf(coder, req)
	}); err != nil {
		logger.Error(err)
	}
	if err := b.RpcX.SubscribeServer(func(coder *rpcx.RpcEncoder, req *rpcx.RpcMsg) []byte {
		logger.Infof("BaseBackEnd SubscribeServer received: %+v", req)
		return b.EventHandlerBroadcast(coder, req)
	}); err != nil {
		logger.Error(err)
	}
	//register the service
	if err := discover.Register(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *BaseBackEnd) BeforeShutdown() {
	//unregister the service
	if err := discover.UnRegister(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *BaseBackEnd) Shutdown() {
	//shutdown server
	logger.Infoln("stop the backend:", b.ServerId)
}

func (b *BaseBackEnd) GetServer() *treaty.Server {
	return b.Server
}

func (b *BaseBackEnd) RegEventHandlerSelf(handler rpcx.CallbackFunc) { //注册自己事件处理器
	b.EventHandlerSelf = handler
}

func (b *BaseBackEnd) RegEventHandlerBroadcast(handler rpcx.CallbackFunc) { //注册广播事件处理器
	b.EventHandlerBroadcast = handler
}

func (b *BaseBackEnd) SetServerId(serverId string) {
	b.ServerId = serverId
}

func (b *BaseBackEnd) GetServerId() string {
	return b.ServerId
}
