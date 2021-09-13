package backend

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/kungfu/utils"
)

//DatabaseEnd 数据库后端服务
type DatabaseEnd struct {
	ServerId              string
	Server                *treaty.Server
	RpcX                  rpcx.RpcServer
	EventHandlerSelf      rpcx.CallbackFunc //处理自己的事件
	EventHandlerBroadcast rpcx.CallbackFunc //处理广播事件
}

func (b *DatabaseEnd) Init() {
	//find the  server config
	if b.Server = utils.FindServerConfig(config.GetServersConf(), b.GetServerId()); b.Server == nil {
		logger.Fatal("DatabaseEnd can find the server config")
	}
	//init the rpcx
	b.RpcX = rpcx.NewRpcServer(config.GetRpcXConf(), b.Server)
	logger.Info("init the backend:", b.ServerId)
}

func (b *DatabaseEnd) AfterInit() {
	//Subscribe event
	if err := b.RpcX.Subscribe(b.Server, func(req *rpcx.RpcMsg) []byte {
		//logger.Infof("BaseBackEnd Subscribe received: %+v", req)
		return b.EventHandlerSelf(req)
	}); err != nil {
		logger.Error(err)
	}
	//Subscribe event
	if err := b.RpcX.SubscribeJson(b.Server, func(req *rpcx.RpcMsg) []byte {
		//logger.Infof("DatabaseEnd Subscribe received: %+v", req)
		return b.EventHandlerSelf(req)
	}, false); err != nil {
		logger.Error(err)
	}
	if err := b.RpcX.SubscribeDatabase(func(req *rpcx.RpcMsg) []byte {
		logger.Infof("DatabaseEnd SubscribeDatabase received: %+v", req)
		return b.EventHandlerBroadcast(req)
	}); err != nil {
		logger.Error(err)
	}
	//register the service
	if err := discover.Register(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *DatabaseEnd) BeforeShutdown() {
	//unregister the service
	if err := discover.UnRegister(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *DatabaseEnd) Shutdown() {
	//shutdown server
	logger.Info("stop the backend:", b.ServerId)
}

func (b *DatabaseEnd) GetServer() *treaty.Server {
	return b.Server
}

func (b *DatabaseEnd) RegEventHandlerSelf(handler rpcx.CallbackFunc) { //注册自己事件处理器
	b.EventHandlerSelf = handler
}

func (b *DatabaseEnd) RegEventHandlerBroadcast(handler rpcx.CallbackFunc) { //注册广播事件处理器
	b.EventHandlerBroadcast = handler
}

func (b *DatabaseEnd) SetServerId(serverId string) {
	b.ServerId = serverId
}

func (b *DatabaseEnd) GetServerId() string {
	return b.ServerId
}
