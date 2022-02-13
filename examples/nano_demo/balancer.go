package main

import (
	"errors"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"

	"github.com/jqiris/kungfu/balancer"
	"github.com/jqiris/kungfu/launch"
)

type MyBalancer struct {
	balancer.BaseBalancer
}

func (b *MyBalancer) EventHandleSelf(req *rpcx.RpcMsg) []byte {
	logger.Infof("MyBalancer EventHandleSelf received: %+v", req)
	return nil
}

func (b *MyBalancer) EventHandleBroadcast(req *rpcx.RpcMsg) []byte {
	logger.Infof("MyBalancer EventHandleBroadcast received: %+v", req)
	return nil
}

func MyBalancerCreator(s *treaty.Server) (rpcx.ServerEntity, error) {
	if len(s.ServerId) < 1 {
		return nil, errors.New("服务器id不能为空")
	}
	server := &MyBalancer{
		BaseBalancer: balancer.BaseBalancer{
			Server: s,
		},
	}
	server.BaseBalancer.EventJsonSelf = server.EventHandleSelf
	server.BaseBalancer.EventHandlerSelf = server.EventHandleSelf
	server.BaseBalancer.EventHandlerBroadcast = server.EventHandleBroadcast
	return server, nil
}

func init() {
	launch.RegisterCreator(rpcx.Balancer, MyBalancerCreator)
}
