package main

import (
	"fmt"
	"github.com/jqiris/kungfu/rpc"

	"github.com/jqiris/kungfu/balancer"
	"github.com/jqiris/kungfu/launch"
)

type MyBalancer struct {
	balancer.BaseBalancer
}

func (b *MyBalancer) HandleSelfEvent(server rpc.ServerRpc, req *rpc.MsgRpc) []byte {
	fmt.Printf("MyBalancer HandleSelfEvent received: %+v \n", req)
	return nil
}

func (b *MyBalancer) HandleBroadcastEvent(server rpc.ServerRpc, req *rpc.MsgRpc) []byte {
	fmt.Printf("MyBalancer HandleBroadcastEvent received: %+v \n", req)
	return nil
}

func init() {
	srv := &MyBalancer{}
	srv.SetServerId("balancer_1001")
	srv.RegEventHandlerSelf(srv.HandleSelfEvent)
	srv.RegEventHandlerBroadcast(srv.HandleBroadcastEvent)
	launch.RegisterServer(srv)

	srv2 := &MyBalancer{}
	srv2.SetServerId("balancer_1002")
	srv2.RegEventHandlerSelf(srv2.HandleSelfEvent)
	srv2.RegEventHandlerBroadcast(srv2.HandleBroadcastEvent)
	launch.RegisterServer(srv2)
}
