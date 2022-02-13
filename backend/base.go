package backend

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
)

//BaseBackEnd 后盾服务
type BaseBackEnd struct {
	ServerId              string
	Server                *treaty.Server
	RpcX                  rpcx.RpcServer
	EventJsonSelf         rpcx.CallbackFunc //处理自己的json事件
	EventHandlerSelf      rpcx.CallbackFunc //处理自己的事件
	EventHandlerBroadcast rpcx.CallbackFunc //处理广播事件
}

func (b *BaseBackEnd) Init() {
	if b.Server == nil {
		panic("服务配置信息不能为空")
		return
	}
	//赋值id
	b.ServerId = b.Server.ServerId
	//init the rpcx
	b.RpcX = rpcx.NewRpcServer(config.GetRpcXConf(), b.Server)
	logger.Info("init the backend:", b.ServerId)
}

func (b *BaseBackEnd) AfterInit() {
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
	if err := b.RpcX.SubscribeServer(*builder); err != nil {
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
	logger.Info("stop the backend:", b.ServerId)
}

func (b *BaseBackEnd) GetServer() *treaty.Server {
	return b.Server
}

func (b *BaseBackEnd) RegEventJsonSelf(handler rpcx.CallbackFunc) { //注册自己事件处理器
	b.EventJsonSelf = handler
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
