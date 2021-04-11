package tests

import (
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/treaty"
	"testing"
)

func TestEtcdDisCover(t *testing.T) {
	server := treaty.Server{
		ServerId:   1001,
		ServerType: treaty.ServerGate,
		ServerName: "gate",
		ServerHost: "127.0.0.1:1234",
	}
	////reg server
	//err := discover.Register(server)
	//if err != nil {
	//	logger.Error(err)
	//}

	//get sever list
	res := discover.DiscoverServerList()
	logger.Infof("the server list:%+v", res)
	res2 := discover.DiscoverServer(treaty.ServerGate)
	logger.Infof("the server list:%+v", res2)
	//unregister server
	err := discover.UnRegister(server)
	if err != nil {
		logger.Errorf("UnRegister err:%v", err)
	}
	res = discover.DiscoverServerList()
	logger.Infof("the server list:%+v", res)
}
