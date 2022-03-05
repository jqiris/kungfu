package plugin

import (
	"fmt"

	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/treaty"
)

type ServerHttp struct {
	handler rpc.HttpHandler
}

func NewServerHttp(h rpc.HttpHandler) *ServerHttp {
	return &ServerHttp{
		handler: h,
	}
}
func (b *ServerHttp) Init(s *treaty.Server) {
}

func (b *ServerHttp) Run(s *treaty.Server) {
	if b.handler == nil {
		panic("http handler is nil")
	}
	addr := fmt.Sprintf(":%d", s.ClientPort)
	if err := b.handler.Run(addr); err != nil {
		panic(err)
	}
}

func (b *ServerHttp) AfterInit(s *treaty.Server) {
	go b.Run(s)
}

func (b *ServerHttp) BeforeShutdown() {
}

func (b *ServerHttp) Shutdown() {
}
