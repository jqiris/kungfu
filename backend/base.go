package backend

import (
	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
)

//BaseBackEnd
type BaseBackEnd struct {
	ServerId              string
	Server                *treaty.Server
	RpcX                  rpcx.RpcServer
	EventHandlerSelf      func(req []byte) []byte //处理自己的事件
	EventHandlerBroadcast func(req []byte) []byte //处理广播事件
}

func (b *BaseBackEnd) Init() {
	//find the  server config
	if b.Server = helper.FindServerConfig(conf.GetServersConf(), b.GetServerId()); b.Server == nil {
		logger.Fatal("BaseBackEnd can find the server config")
	}
	//init the rpcx
	b.RpcX = rpcx.NewRpcServer(conf.GetRpcxConf())
	logger.Infoln("init the backend:", b.ServerId)
}

func (b *BaseBackEnd) AfterInit() {
	//Subscribe event
	if err := b.RpcX.Subscribe(b.Server, func(req []byte) []byte {
		logger.Infof("BaseBackEnd Subscribe received: %+v", req)
		return b.EventHandlerSelf(req)
	}); err != nil {
		logger.Error(err)
	}
	if err := b.RpcX.SubscribeServer(func(req []byte) []byte {
		logger.Infof("BaseBackEnd SubscribeServer received: %+v", req)
		return b.EventHandlerBroadcast(req)
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

func (b *BaseBackEnd) RegEventHandlerSelf(handler func(req []byte) []byte) { //注册自己事件处理器
	b.EventHandlerSelf = handler
}

func (b *BaseBackEnd) RegEventHandlerBroadcast(handler func(req []byte) []byte) { //注册广播事件处理器
	b.EventHandlerBroadcast = handler
}

func (b *BaseBackEnd) SetServerId(serverId string) {
	b.ServerId = serverId
}

func (b *BaseBackEnd) GetServerId() string {
	return b.ServerId
}
