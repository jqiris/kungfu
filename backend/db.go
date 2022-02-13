package backend

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpc"
	"github.com/jqiris/kungfu/treaty"
)

//DatabaseEnd 数据库后端服务
type DatabaseEnd struct {
	ServerId              string
	Server                *treaty.Server
	Rpc                   rpc.ServerRpc
	EventJsonSelf         rpc.CallbackFunc //处理自己的json事件
	EventHandlerSelf      rpc.CallbackFunc //处理自己的事件
	EventHandlerBroadcast rpc.CallbackFunc //处理广播事件
}

func (b *DatabaseEnd) Init() {
	if b.Server == nil {
		panic("服务配置信息不能为空")
		return
	}
	//赋值id
	b.ServerId = b.Server.ServerId
	//init the rpc
	b.Rpc = rpc.NewRpcServer(config.GetRpcConf(), b.Server)
	logger.Info("init the backend:", b.ServerId)
}

func (b *DatabaseEnd) AfterInit() {
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
	builder := rpc.NewSubscriberRpc(b.Server).SetCodeType(rpc.CodeTypeProto).SetCallback(func(req *rpc.MsgRpc) []byte {
		return b.EventHandlerSelf(req)
	})
	//Subscribe event
	if err := b.Rpc.Subscribe(builder.Build()); err != nil {
		logger.Error(err)
	}
	builder = builder.SetSuffix("json").SetCodeType(rpc.CodeTypeJson).SetCallback(func(req *rpc.MsgRpc) []byte {
		return b.EventJsonSelf(req)
	})
	//Subscribe event
	if err := b.Rpc.Subscribe(builder.Build()); err != nil {
		logger.Error(err)
	}
	builder = builder.SetSuffix(rpc.DefaultSuffix).SetCodeType(rpc.CodeTypeProto).SetCallback(func(req *rpc.MsgRpc) []byte {
		return b.EventHandlerBroadcast(req)
	})
	if err := b.Rpc.SubscribeDatabase(builder.Build()); err != nil {
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

func (b *DatabaseEnd) SetServerId(serverId string) {
	b.ServerId = serverId
}

func (b *DatabaseEnd) GetServerId() string {
	return b.ServerId
}
