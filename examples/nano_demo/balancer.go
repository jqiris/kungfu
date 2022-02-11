package main

import (
	"errors"
	"fmt"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"

	"github.com/jqiris/kungfu/balancer"
	"github.com/jqiris/kungfu/launch"
)

type MyBalancer struct {
	balancer.BaseBalancer
}

func EventHandleSelf(req *rpcx.RpcMsg) []byte {
	fmt.Printf("MyBalancer EventHandleSelf received: %+v \n", req)
	return nil
}

func EventHandleBroadcast(req *rpcx.RpcMsg) []byte {
	fmt.Printf("MyBalancer EventHandleBroadcast received: %+v \n", req)
	return nil
}

func MyBalancerCreator(s *treaty.Server) (rpcx.ServerEntity, error) {
	if len(s.ServerId) < 1 {
		return nil, errors.New("服务器id不能为空")
	}
	server := &MyBalancer{
		BaseBalancer: balancer.BaseBalancer{
			Server:                s,
			EventJsonSelf:         EventHandleSelf,
			EventHandlerSelf:      EventHandleSelf,
			EventHandlerBroadcast: EventHandleBroadcast,
		},
	}
	return server, nil
}

func init() {
	launch.RegisterCreator(rpcx.Balancer, MyBalancerCreator)
}
