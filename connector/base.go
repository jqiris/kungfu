package connector

import (
	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/zinx/utils"
	"github.com/jqiris/zinx/ziface"
	"github.com/jqiris/zinx/znet"
)

type BaseConnector struct {
	Server                *treaty.Server
	Rpcx                  rpcx.RpcConnector
	ClientServer          ziface.IServer
	ClientCoder           coder.Coder
	ConnectorConf         *utils.GlobalObj
	EventHandlerSelf      func(req []byte) []byte //处理自己的事件
	EventHandlerBroadcast func(req []byte) []byte //处理广播事件
}

func (b *BaseConnector) Init() {
	//run the front server
	b.ClientServer = znet.NewServer(*b.ConnectorConf)
	go b.ClientServer.Serve()
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
