package connector

import (
	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/zinx/utils"
	"github.com/jqiris/zinx/ziface"
	"github.com/jqiris/zinx/znet"
)

type BaseConnector struct {
	ServerId              int32
	Server                *treaty.Server
	Rpcx                  rpcx.RpcConnector
	ClientServer          ziface.IServer
	ClientCoder           coder.Coder
	ConnectorConf         *utils.GlobalObj
	EventHandlerSelf      func(req []byte) []byte //处理自己的事件
	EventHandlerBroadcast func(req []byte) []byte //处理广播事件
}

func (b *BaseConnector) Init() {
	//find the  server config
	if b.ConnectorConf = helper.FindConnectorConfig(conf.GetConnectorConf(), b.GetServerId()); b.ConnectorConf == nil {
		logger.Fatal("BaseConnector can find the server config")
	} else {
		b.Server = &treaty.Server{
			ServerId:   b.ConnectorConf.ServerId,
			ServerType: treaty.ServerType(b.ConnectorConf.ServerType),
			ServerName: b.ConnectorConf.ServerName,
			ServerIp:   b.ConnectorConf.ServerIp,
			ClientPort: int32(b.ConnectorConf.ClientPort),
		}
	}
	//init the rpcx
	b.Rpcx = rpcx.NewRpcConnector(conf.GetRpcxConf())
	//run the front server
	b.ClientServer = znet.NewServer(*b.ConnectorConf)
	go b.ClientServer.Serve()

	logger.Infoln("init the connector:", b.ServerId)
}

func (b *BaseConnector) AfterInit() {
	//Subscribe event
	if err := b.Rpcx.Subscribe(b.Server, func(req []byte) []byte {
		logger.Infof("BaseConnector Subscribe received: %+v", req)
		return b.EventHandlerSelf(req)
	}); err != nil {
		logger.Error(err)
	}
	if err := b.Rpcx.SubscribeConnector(func(req []byte) []byte {
		logger.Infof("BaseConnector SubscribeConnector received: %+v", req)
		return b.EventHandlerBroadcast(req)
	}); err != nil {
		logger.Error(err)
	}
	//register the service
	if err := discover.Register(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *BaseConnector) BeforeShutdown() {
	//unregister the service
	if err := discover.UnRegister(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *BaseConnector) Shutdown() {
	//stop the server
	if b.ClientServer != nil {
		b.ClientServer.Stop()
	}
}

func (b *BaseConnector) GetServer() *treaty.Server {
	return b.Server
}

func (b *BaseConnector) RegEventHandlerSelf(handler func(req []byte) []byte) { //注册自己事件处理器
	b.EventHandlerSelf = handler
}

func (b *BaseConnector) RegEventHandlerBroadcast(handler func(req []byte) []byte) { //注册广播事件处理器
	b.EventHandlerBroadcast = handler
}
func (b *BaseConnector) SetServerId(serverId int32) {
	b.ServerId = serverId
}

func (b *BaseConnector) GetServerId() int32 {
	return b.ServerId
}
