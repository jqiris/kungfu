package launch

import (
	"sync"

	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/treaty"
	"github.com/sirupsen/logrus"
)

//服务器集群管理
var (
	logger     = logrus.WithField("package", "launch")
	servers    map[int32]treaty.ServerEntity
	launched   map[int32]treaty.ServerEntity
	launchLock *sync.Mutex
)

func init() {
	servers = make(map[int32]treaty.ServerEntity)
	launched = make(map[int32]treaty.ServerEntity)
	launchLock = new(sync.Mutex)
}

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
	//run servers
	for _, server := range servers {
		if conf.IsInLauch(server.GetServerId()) {
			go func(srv treaty.ServerEntity) {
				srv.Init()
				srv.AfterInit()
				launchLock.Lock()
				defer launchLock.Unlock()
				launched[srv.GetServerId()] = srv
			}(server)
		}
	}

	<-done
	//stop servers
	for _, server := range launched {
		go func(srv treaty.ServerEntity) {
			srv.BeforeShutdown()
			srv.Shutdown()
		}(server)
	}
}
