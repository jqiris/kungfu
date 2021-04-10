package tests

import (
	"github.com/jqiris/kungfu/common"
	"github.com/jqiris/kungfu/discover"
	"testing"
)

func TestEtcdDisCover(t *testing.T) {
	server := common.Server{
		ServerId:   1001,
		ServerType: common.ServerGate,
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
	res2 := discover.DiscoverServer(common.ServerGate)
	logger.Infof("the server list:%+v", res2)
	//unregister server
	err := discover.UnRegister(server)
	if err != nil {
		logger.Errorf("UnRegister err:%v", err)
	}
	res = discover.DiscoverServerList()
	logger.Infof("the server list:%+v", res)
}
