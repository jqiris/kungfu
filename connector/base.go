package connector

import (
	"fmt"
	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
	"net"
	"net/http"
)

type BaseConnector struct {
	Server       *treaty.Server
	Rpcx         rpcx.RpcBalancer
	ClientServer *http.Server
	ClientCoder  coder.Coder
}

func (b *BaseConnector) Init() {
	//run the front server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", b.Server.ClientPort))
	if err != nil {
		logger.Fatal(err)
		return
	}
	go b.RunClientServer(listener)
}

func (b *BaseConnector) RunClientServer(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error(err)
			return
		}
		go b.DealClientConnection(conn)
	}
}
func (b *BaseConnector) DealClientConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)

	}
}

func (b *BaseConnector) AfterInit() {
	//register the server
}

func (b *BaseConnector) BeforeShutdown() {
	panic("implement me")
}

func (b *BaseConnector) Shutdown() {
	panic("implement me")
}
