package plugin

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/rpc"
	"net/http"
	"sync"
)

var (
	once = sync.Once{}
)

type ServerSocket struct {
	sc *socketio.Server
	ns string
}

func NewServerSocket(ns string, opts *engineio.Options) *ServerSocket {
	return &ServerSocket{
		ns: ns,
		sc: socketio.NewServer(opts),
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
	go b.sc.Serve()
	defer b.sc.Close()
	once.Do(func() {
		http.Handle("/socket.io/", b.sc)
	})
	http.Handle("/socket.io/", b.sc)
	logger.Infof("socket server start at:%v", s.Server.ClientPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", s.Server.ClientPort), nil); err != nil {
		logger.Fatal(err)
	}
}

func (b *ServerSocket) Init(s *rpc.ServerBase) {
}

func (b *ServerSocket) AfterInit(s *rpc.ServerBase) {
	go b.Run(s)
}

func (b *ServerSocket) BeforeShutdown(s *rpc.ServerBase) {
}

func (b *ServerSocket) Shutdown(s *rpc.ServerBase) {
}
