package examples

import (
	"fmt"

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

func (b *MyConnector) EventHandleBroasdcast(req []byte) []byte {
	fmt.Printf("MyConnector EventHandleBroasdcast received: %+v \n", string(req))
	return nil
}

func init() {
	srv := &MyConnector{}
	srv.SetServerId(2001)
	srv.RegEventHandlerSelf(srv.EventHandlerSelf)
	srv.RegEventHandlerBroadcast(srv.EventHandleBroasdcast)
	launch.RegisterServer(srv)
}
