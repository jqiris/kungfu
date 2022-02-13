package main

import (
	"github.com/jqiris/kungfu/base"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpc"
	"github.com/jqiris/kungfu/treaty"

	"github.com/jqiris/kungfu/launch"
)

type MyBalancer struct {
	*base.ServerBalancer
}

func (b *MyBalancer) HandleSelfEvent(req *rpc.MsgRpc) []byte {
	logger.Infof("MyBalancer HandleSelfEvent received: %+v", req)
	return nil
}

func (b *MyBalancer) HandleBroadcastEvent(req *rpc.MsgRpc) []byte {
	logger.Infof("MyBalancer HandleBroadcastEvent received: %+v", req)
	return nil
}

func MyBalancerCreator(s *treaty.Server) (rpc.ServerEntity, error) {
	server := &MyBalancer{
		ServerBalancer: base.NewServerBalancer(s),
	}
	server.SelfEventHandler = server.HandleSelfEvent
	server.BroadcastEventHandler = server.HandleBroadcastEvent
	return server, nil
}

func init() {
	launch.RegisterCreator(rpc.Balancer, MyBalancerCreator)
}
