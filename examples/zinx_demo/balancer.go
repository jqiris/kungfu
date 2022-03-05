package main

import (
	"github.com/jqiris/kungfu/v2/base/plugin"
	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/treaty"

	"github.com/jqiris/kungfu/v2/base"
	"github.com/jqiris/kungfu/v2/launch"
)

type MyBalancer struct {
	*base.ServerBase
}

func MyBalancerCreator(s *treaty.Server) (rpc.ServerEntity, error) {
	server := &MyBalancer{
		ServerBase: base.NewServerBase(s),
	}
	plug := plugin.NewServerBalancer()
	server.AddPlugin(plug)
	return server, nil
}

func init() {
	launch.RegisterCreator("balancer", MyBalancerCreator)
}
