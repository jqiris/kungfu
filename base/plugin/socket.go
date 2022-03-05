package plugin

import (
	"fmt"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/treaty"
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

func (b *ServerSocket) Run(s *treaty.Server) {
	go b.sc.Serve()
	defer b.sc.Close()

	http.Handle("/socket.io/", b.sc)
	logger.Infof("socket server start at:%v", s.ClientPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", s.ClientPort), nil); err != nil {
		logger.Fatal(err)
	}
}

func (b *ServerSocket) Init(s *treaty.Server) {
}

func (b *ServerSocket) AfterInit(s *treaty.Server) {
	go b.Run(s)
}

func (b *ServerSocket) BeforeShutdown() {
}

func (b *ServerSocket) Shutdown() {
}
