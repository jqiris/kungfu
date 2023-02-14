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
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/utils"
)

type ServerSocket struct {
	sc  *socketio.Server
	ns  string
	ssl config.SslConf
}

func NewServerSocket(ns string, ssl config.SslConf, opts *engineio.Options) *ServerSocket {
	return &ServerSocket{
		ns:  ns,
		ssl: ssl,
		sc:  socketio.NewServer(opts),
	}
}

func (b *ServerSocket) OnConnect(f func(socketio.Conn) error) {
	b.sc.OnConnect(b.ns, f)
}
func (b *ServerSocket) OnDisconnect(f func(socketio.Conn, string)) {
	b.sc.OnDisconnect(b.ns, f)
}

func (b *ServerSocket) OnEvent(event string, f interface{}) {
	b.sc.OnEvent(b.ns, event, f)
}
func (b *ServerSocket) OnError(f func(socketio.Conn, error)) {
	b.sc.OnError(b.ns, f)
}

func (b *ServerSocket) Run(s *rpc.ServerBase) {
	go utils.SafeRun(func() {
		defer b.sc.Close()
		b.sc.Serve()
	})
	scAddr := "/" + s.Server.ServerId + "/"
	http.Handle(scAddr, b.sc)
	logger.Infof("socket server start at:%v", s.Server.ClientPort)
	if b.ssl.PowerOn {
		if err := http.ListenAndServeTLS(fmt.Sprintf(":%v", s.Server.ClientPort), b.ssl.CertFile, b.ssl.KeyFile, nil); err != nil {
			logger.Fatal(err)
		}
	} else {
		if err := http.ListenAndServe(fmt.Sprintf(":%v", s.Server.ClientPort), nil); err != nil {
			logger.Fatal(err)
		}
	}

}

func (b *ServerSocket) Init(s *rpc.ServerBase) {
}

func (b *ServerSocket) AfterInit(s *rpc.ServerBase) {
	go utils.SafeRun(func() {
		b.Run(s)
	})
}

func (b *ServerSocket) BeforeShutdown(s *rpc.ServerBase) {
}

func (b *ServerSocket) Shutdown(s *rpc.ServerBase) {
}
