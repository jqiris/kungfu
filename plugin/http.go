package plugin

import (
	"fmt"

	"github.com/jqiris/kungfu/v2/rpc"
)

type ServerHttp struct {
	handler rpc.HttpHandler
}

func NewServerHttp(h rpc.HttpHandler) *ServerHttp {
	return &ServerHttp{
		handler: h,
	}
}
func (b *ServerHttp) Init(s *rpc.ServerBase) {
}

func (b *ServerHttp) Run(s *rpc.ServerBase) {
	if b.handler == nil {
		panic("http handler is nil")
	}
	addr := fmt.Sprintf(":%d", s.Server.ClientPort)
	if err := b.handler.Run(addr); err != nil {
		panic(err)
	}
}

func (b *ServerHttp) AfterInit(s *rpc.ServerBase) {
	go b.Run(s)
}

func (b *ServerHttp) BeforeShutdown(s *rpc.ServerBase) {
}

func (b *ServerHttp) Shutdown(s *rpc.ServerBase) {
}
