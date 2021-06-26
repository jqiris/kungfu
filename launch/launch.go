package launch

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/sirupsen/logrus"
)

//服务器集群管理
var (
	logger   = logrus.WithField("package", "launch")
	servers  map[string]rpcx.ServerEntity
	launched map[string]rpcx.ServerEntity
)

func init() {
	servers = make(map[string]rpcx.ServerEntity)
	launched = make(map[string]rpcx.ServerEntity)
}

func RegisterServer(server rpcx.ServerEntity) {
	if _, ok := servers[server.GetServerId()]; !ok {
		servers[server.GetServerId()] = server
	} else {
		logger.Errorf("RegisterServer duplicate, error: %+v", server)
	}
}

func UnRegisterServer(server rpcx.ServerEntity) {
	delete(servers, server.GetServerId())
}

//Startup 启动服务器
func Startup() {
	//run servers
	for _, server := range servers {
		if config.IsInLaunch(server.GetServerId()) {
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
