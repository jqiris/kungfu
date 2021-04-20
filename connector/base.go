package connector

import (
	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/zinx/ziface"
	"github.com/jqiris/zinx/znet"
)

type BaseConnector struct {
	Server       *treaty.Server
	Rpcx         rpcx.RpcBalancer
	ClientServer ziface.IServer
	ClientCoder  coder.Coder
}

func (b *BaseConnector) Init() {
	//run the front server
	b.ClientServer = znet.NewServer()
	b.ClientServer.Serve()
}

func (b *BaseConnector) AfterInit() {
	//register the server
}

func (b *BaseConnector) BeforeShutdown() {
	panic("implement me")
}

func (b *BaseConnector) Shutdown() {
	panic("implement me")
}
