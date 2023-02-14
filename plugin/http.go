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

package plugin

import (
	"fmt"

	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/utils"
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
	addr := fmt.Sprintf(":%d", s.Server.ClientPort)
	if err := b.handler.Run(addr); err != nil {
		panic(err)
	}
}

func (b *ServerHttp) AfterInit(s *rpc.ServerBase) {
	if b.handler == nil {
		panic("http handler is nil")
	}
	go utils.SafeRun(func() {
		b.Run(s)
	})
}

func (b *ServerHttp) BeforeShutdown(s *rpc.ServerBase) {
}

func (b *ServerHttp) Shutdown(s *rpc.ServerBase) {
}
