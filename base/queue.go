package base

import (
	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/treaty"
)

type ServerQueue struct {
	*ServerBase
	queue  string
	suffix string
}

func NewServerQueue(s *treaty.Server, queue string) *ServerQueue {
	return &ServerQueue{
		ServerBase: NewServerBase(s),
		queue:      queue,
	}
}

func (s *ServerQueue) AfterInit() {
	s.ServerBase.AfterInit()
	//订阅queue消息
	b := s.SubBuilder.Build()
	b = b.SetQueue(s.queue).Build()
	//self queue proto event
	if err := s.Rpc.QueueSubscribe(b.SetCodeType(rpc.CodeTypeProto).SetCallback(s.SelfEventHandler).Build()); err != nil {
		panic(err)
	}
	//self queue json event
	if err := s.Rpc.QueueSubscribe(b.SetSuffix(rpc.JsonSuffix).SetCodeType(rpc.CodeTypeJson).SetCallback(s.SelfEventHandler).Build()); err != nil {
		panic(err)
	}
}
