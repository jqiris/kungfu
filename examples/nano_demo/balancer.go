package main

import (
	"errors"
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
	if len(s.ServerId) < 1 {
		return nil, errors.New("服务器id不能为空")
	}
	server := &MyBalancer{
		ServerBalancer: base.NewServerBalancer(s),
	}
	return server, nil
}

func init() {
	launch.RegisterCreator(rpc.Balancer, MyBalancerCreator)
}
