package main

import (
	"github.com/jqiris/kungfu/v2/plugin"
	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/treaty"

	"github.com/jqiris/kungfu/v2/launch"
)

type MyBalancer struct {
	*rpc.ServerBase
}

func MyBalancerCreator(s *treaty.Server) (rpc.ServerEntity, error) {
	server := &MyBalancer{
		ServerBase: rpc.NewServerBase(s),
	}
	plug := plugin.NewServerBalancer()
	server.AddPlugin(plug)
	return server, nil
}

func init() {
	launch.RegisterCreator(rpc.Balancer, MyBalancerCreator)
}
