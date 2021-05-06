package examples

import (
	"fmt"
	"github.com/jqiris/kungfu/examples/handler"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/zinx/ziface"

	"github.com/jqiris/kungfu/connector"
	"github.com/jqiris/kungfu/launch"
)

type MyConnector struct {
	connector.BaseConnector
}

func (b *MyConnector) EventHandleSelf(req []byte) []byte {
	fmt.Printf("MyConnector EventHandleSelf received: %+v \n", string(req))
	return nil
}

func (b *MyConnector) EventHandleBroadcast(req []byte) []byte {
	fmt.Printf("MyConnector EventHandleBroadcast received: %+v \n", string(req))
	return nil
}

func init() {
	routers := map[uint32]ziface.IRouter{
		uint32(treaty.MsgId_Msg_Login_Request): &handler.LogingHandler{},
	}
	srv := &MyConnector{}
	srv.SetServerId("connector_2001")
	srv.RegEventHandlerSelf(srv.EventHandlerSelf)
	srv.RegEventHandlerBroadcast(srv.EventHandleBroadcast)
	srv.RegRouters(routers)
	launch.RegisterServer(srv)

	srv2 := &MyConnector{}
	srv2.SetServerId("connector_2002")
	srv2.RegEventHandlerSelf(srv2.EventHandlerSelf)
	srv2.RegEventHandlerBroadcast(srv2.EventHandleBroadcast)
	srv2.RegRouters(routers)
	launch.RegisterServer(srv2)
}
