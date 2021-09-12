package main

import (
	"fmt"
	"github.com/jqiris/kungfu/rpcx"

	"github.com/jqiris/kungfu/balancer"
	"github.com/jqiris/kungfu/launch"
)

type MyBalancer struct {
	balancer.BaseBalancer
}

func (b *MyBalancer) EventHandleSelf(req *rpcx.RpcMsg) []byte {
	fmt.Printf("MyBalancer EventHandleSelf received: %+v \n", req)
	return nil
}

func (b *MyBalancer) EventHandleBroadcast(req *rpcx.RpcMsg) []byte {
	fmt.Printf("MyBalancer EventHandleBroadcast received: %+v \n", req)
	return nil
}

func init() {
	srv := &MyBalancer{}
	srv.SetServerId("balancer_1001")
	srv.RegEventHandlerSelf(srv.EventHandleSelf)
	srv.RegEventHandlerBroadcast(srv.EventHandleBroadcast)
	launch.RegisterServer(srv)

	srv2 := &MyBalancer{}
	srv2.SetServerId("balancer_1002")
	srv2.RegEventHandlerSelf(srv2.EventHandleSelf)
	srv2.RegEventHandlerBroadcast(srv2.EventHandleBroadcast)
	launch.RegisterServer(srv2)
}
