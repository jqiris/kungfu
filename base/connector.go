package base

import (
	"github.com/jqiris/kungfu/v2/tcpface"
	"github.com/jqiris/kungfu/v2/tcpserver"
	"github.com/jqiris/kungfu/v2/treaty"
)

type ServerConnector struct {
	*ServerBase
	ClientServer tcpface.IServer         //client server
	RouteHandler func(s tcpface.IServer) //注册路由
}

func NewServerConnector(s *treaty.Server) *ServerConnector {
	return &ServerConnector{
		ServerBase: NewServerBase(s),
	}
}

func (s *ServerConnector) Init() {
	s.ServerBase.Init()
	if s.RouteHandler == nil {
		panic("连接器路由配置信息不能为空")
	}
	//run the front server
	s.ClientServer = tcpserver.NewServer(s.Server)
	s.RouteHandler(s.ClientServer)
	go s.ClientServer.Serve()
}

func (s *ServerConnector) Shutdown() {
	s.ServerBase.Shutdown()
	//stop the server
	if s.ClientServer != nil {
		s.ClientServer.Stop()
	}
}
