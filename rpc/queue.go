/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package rpc

import (
	"github.com/jqiris/kungfu/v2/treaty"
)

// ServerQueue base queue
type ServerQueue struct {
	*ServerBase
	queue string
}

func NewServerQueue(queue string, s *treaty.Server, options ...Option) *ServerQueue {
	return &ServerQueue{
		ServerBase: NewServerBase(s, options...),
		queue:      queue,
	}
}

func (s *ServerQueue) AfterInit() {
	s.ServerBase.AfterInit()
	//订阅queue消息
	b := s.SubBuilder.Build()
	b = b.SetQueue(s.queue).Build()
	//self queue proto event
	if err := s.Rpc.QueueSubscribe(b.SetCodeType(CodeTypeProto).SetCallback(s.selfEventHandler).Build()); err != nil {
		panic(err)
	}
	//self queue json event
	if err := s.Rpc.QueueSubscribe(b.SetSuffix(JsonSuffix).SetCodeType(CodeTypeJson).SetCallback(s.selfEventHandler).Build()); err != nil {
		panic(err)
	}
}
