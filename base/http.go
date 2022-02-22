package base

import (
	"fmt"

	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/treaty"
)

type ServerHttp struct {
	*ServerBase
	handler rpc.HttpHandler
}

func NewServerHttp(s *treaty.Server, h rpc.HttpHandler) *ServerHttp {
	return &ServerHttp{
		ServerBase: NewServerBase(s),
		handler:    h,
	}
}

func (s *ServerHttp) Init() {
	s.ServerBase.Init()
}

func (s *ServerHttp) Run() {
	if s.handler == nil {
		panic("http handler is nil")
	}
	addr := fmt.Sprintf(":%d", s.ServerBase.Server.ClientPort)
	if err := s.handler.Run(addr); err != nil {
		panic(err)
	}
}

func (s *ServerHttp) AfterInit() {
	s.ServerBase.AfterInit()
	go s.Run()
}
