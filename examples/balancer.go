package examples

import (
	"fmt"

	"github.com/jqiris/kungfu/balancer"
	"github.com/jqiris/kungfu/launch"
)

type MyBalancer struct {
	balancer.BaseBalancer
}

func (b *MyBalancer) EventHandleSelf(req []byte) []byte {
	fmt.Printf("MyBalancer EventHandleSelf received: %+v \n", string(req))
	return nil
}

func (b *MyBalancer) EventHandleBroasdcast(req []byte) []byte {
	fmt.Printf("MyBalancer EventHandleBroasdcast received: %+v \n", string(req))
	return nil
}

func init() {
	srv := &MyBalancer{}
	srv.SetServerId("balancer_1001")
	srv.RegEventHandlerSelf(srv.EventHandlerSelf)
	srv.RegEventHandlerBroadcast(srv.EventHandleBroasdcast)
	launch.RegisterServer(srv)

	// srv2 := &MyBalancer{}
	// srv2.SetServerId("balancer_1002")
	// srv2.RegEventHandlerSelf(srv2.EventHandlerSelf)
	// srv2.RegEventHandlerBroadcast(srv2.EventHandleBroasdcast)
	// launch.RegisterServer(srv2)
}
