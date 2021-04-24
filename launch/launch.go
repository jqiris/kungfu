package launch

import (
	"github.com/jqiris/kungfu/treaty"
	"github.com/sirupsen/logrus"
)

//服务器集群管理
var (
	logger  = logrus.WithField("package", "launch")
	servers map[int32]treaty.ServerEntity
)

func RegisterServer(server treaty.ServerEntity) {
	if _, ok := servers[server.GetServerId()]; !ok {
		servers[server.GetServerId()] = server
	} else {
		logger.Errorf("RegisterServer duplicate, error: %+v", server)
	}
}

func UnRegisterServer(server treaty.ServerEntity) {
	delete(servers, server.GetServerId())
}

func LaunchServers(done chan struct{}) {
	//init servers
	for _, server := range servers {
		go server.Init()
	}

	//after init servers
	for _, server := range servers {
		go server.AfterInit()
	}
	<-done
	//server stop
	for _, server := range servers {
		go server.BeforeShutdown()
	}

	for _, server := range servers {
		go server.Shutdown()
	}
}
