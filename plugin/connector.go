package plugin

import (
	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/tcpface"
	"github.com/jqiris/kungfu/v2/tcpserver"
	"github.com/jqiris/kungfu/v2/utils"
)

type ServerConnector struct {
	ClientServer tcpface.IServer         //client server
	RouteHandler func(s tcpface.IServer) //注册路由
}

func NewServerConnector() *ServerConnector {
	return &ServerConnector{}
}

func (b *ServerConnector) Init(s *rpc.ServerBase) {
	if b.RouteHandler == nil {
		panic("连接器路由配置信息不能为空")
	}
	//run the front server
	go utils.SafeRun(func() {
		b.Run(s)
	})
}
func (b *ServerConnector) Run(s *rpc.ServerBase) {
	//run the front server
	b.ClientServer = tcpserver.NewServer(s.Server)
	b.RouteHandler(b.ClientServer)
	b.ClientServer.Serve()
}

func (b *ServerConnector) AfterInit(s *rpc.ServerBase) {
}

func (b *ServerConnector) BeforeShutdown(s *rpc.ServerBase) {
}

func (b *ServerConnector) Shutdown(s *rpc.ServerBase) {
	//stop the server
	if b.ClientServer != nil {
		b.ClientServer.Stop()
	}
}
