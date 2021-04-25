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
	conf := server.GetServer()
	if _, ok := servers[conf.ServerId]; !ok {
		servers[conf.ServerId] = server
	} else {
		logger.Errorf("RegisterServer duplicate, error: %+v", server)
	}
}

func UnRegisterServer(server treaty.ServerEntity) {
	conf := server.GetServer()
	delete(servers, conf.ServerId)
}

func LaunchServers(done chan struct{}) {
	//run servers
	for _, server := range servers {
		go func(srv treaty.ServerEntity) {
			srv.Init()
			srv.AfterInit()
		}(server)
	}

	<-done
	//stop servers
	for _, server := range servers {
		go func(srv treaty.ServerEntity) {
			srv.BeforeShutdown()
			srv.Shutdown()
		}(server)
	}
}
