package plugin

import (
	"github.com/jqiris/kungfu/v2/tcpface"
	"github.com/jqiris/kungfu/v2/tcpserver"
	"github.com/jqiris/kungfu/v2/treaty"
)

type ServerConnector struct {
	ClientServer tcpface.IServer         //client server
	RouteHandler func(s tcpface.IServer) //注册路由
}

func NewServerConnector() *ServerConnector {
	return &ServerConnector{}
}

func (b *ServerConnector) Init(s *treaty.Server) {
	if b.RouteHandler == nil {
		panic("连接器路由配置信息不能为空")
	}
	//run the front server
	b.ClientServer = tcpserver.NewServer(s)
	b.RouteHandler(b.ClientServer)
	go b.ClientServer.Serve()
}

func (b *ServerConnector) AfterInit(s *treaty.Server) {
}

func (b *ServerConnector) BeforeShutdown() {
}

func (b *ServerConnector) Shutdown() {
	//stop the server
	if b.ClientServer != nil {
		b.ClientServer.Stop()
	}
}
