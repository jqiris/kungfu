package connector

import (
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
)

type NanoConnector struct {
	ServerId              string
	Server                *treaty.Server
	RpcX                  rpcx.RpcConnector
	EventHandlerSelf      func(req []byte) []byte //处理自己的事件
	EventHandlerBroadcast func(req []byte) []byte //处理广播事件
}

func (c *NanoConnector) Init() {
	panic("not implemented") // TODO: Implement
}

func (c *NanoConnector) AfterInit() {
	panic("not implemented") // TODO: Implement
}

func (c *NanoConnector) BeforeShutdown() {
	panic("not implemented") // TODO: Implement
}

func (c *NanoConnector) Shutdown() {
	panic("not implemented") // TODO: Implement
}

func (c *NanoConnector) GetServer() *treaty.Server {
	panic("not implemented") // TODO: Implement
}

func (c *NanoConnector) RegEventHandlerSelf(handler func(req []byte) []byte) {
	panic("not implemented") // TODO: Implement
}

func (c *NanoConnector) RegEventHandlerBroadcast(handler func(req []byte) []byte) {
	panic("not implemented") // TODO: Implement
}

func (c *NanoConnector) SetServerId(serverId string) {
	panic("not implemented") // TODO: Implement
}

func (c *NanoConnector) GetServerId() string {
	panic("not implemented") // TODO: Implement
}
