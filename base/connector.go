package base

import (
	"github.com/jqiris/kungfu/rpc"
	"github.com/jqiris/kungfu/tcpface"
	"github.com/jqiris/kungfu/tcpserver"
	"github.com/jqiris/kungfu/treaty"
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

func (s *ServerConnector) AfterInit() {
	s.ServerBase.AfterInit()
	//sub self json event
	b := s.SubBuilder.Build()
	if err := s.Rpc.Subscribe(b.SetSuffix(rpc.JsonSuffix).SetCodeType(rpc.CodeTypeJson).SetCallback(s.HandleSelfEvent).Build()); err != nil {
		panic(err)
	}
}

func (s *ServerConnector) Shutdown() {
	s.ServerBase.Shutdown()
	//stop the server
	if s.ClientServer != nil {
		s.ClientServer.Stop()
	}
}
