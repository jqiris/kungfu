package launch

import (
	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/treaty"
	"github.com/sirupsen/logrus"
)

//服务器集群管理
var (
	logger   = logrus.WithField("package", "launch")
	servers  map[string]treaty.ServerEntity
	launched map[string]treaty.ServerEntity
)

func init() {
	servers = make(map[string]treaty.ServerEntity)
	launched = make(map[string]treaty.ServerEntity)
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

//Startup 启动服务器
func Startup() {
	//run servers
	for _, server := range servers {
		if conf.IsInLaunch(server.GetServerId()) {
			server.Init()
			server.AfterInit()
			launched[server.GetServerId()] = server
		}
	}
}

//Shutdown 关闭服务器
func Shutdown() {
	//stop servers
	for _, server := range launched {
		server.BeforeShutdown()
		server.Shutdown()
	}
}
