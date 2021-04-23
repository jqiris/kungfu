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
	Server        *treaty.Server
	Rpcx          rpcx.RpcConnector
	ClientServer  ziface.IServer
	ClientCoder   coder.Coder
	ConnectorConf utils.GlobalObj
}

func (b *BaseConnector) Init() {
	//run the front server
	b.ClientServer = znet.NewServer(b.ConnectorConf)
	go b.ClientServer.Serve()
}

func (b *BaseConnector) AfterInit() {
	//Subscribe event
	if err := b.Rpcx.Subscribe(b.Server, func(req []byte) []byte {
		logger.Infof("BaseConnector Subscribe received: %+v", req)
		return nil
	}); err != nil {
		logger.Error(err)
	}
	if err := b.Rpcx.SubscribeConnector(func(req []byte) []byte {
		logger.Infof("BaseConnector SubscribeConnector received: %+v", req)
		return nil
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
