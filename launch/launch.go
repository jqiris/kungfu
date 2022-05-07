package launch

import (
	"sort"

	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/treaty"
)

//服务器集群管理
var (
	creators  map[string]rpc.ServerCreator
	launched  map[string]rpc.ServerEntity
	launchArr []*treaty.Server
)

func init() {
	creators = make(map[string]rpc.ServerCreator)
	launched = make(map[string]rpc.ServerEntity)
	launchArr = make([]*treaty.Server, 0)
}

func RegisterCreator(typ string, creator rpc.ServerCreator) {
	if _, ok := creators[typ]; !ok {
		creators[typ] = creator
	} else {
		logger.Fatalf("RegisterCreator duplicate, type:%v", typ)
	}
}

//Startup 启动服务器
func Startup() {
	//run servers
	servers := config.GetServersConf()
	for _, cfg := range servers {
		if cfg.IsLaunch {
			launchArr = append(launchArr, cfg)
		}
	}
	sort.Slice(launchArr, func(i, j int) bool {
		return launchArr[i].LaunchWeight < launchArr[j].LaunchWeight
	})
	for _, cfg := range launchArr {
		creator := creators[cfg.ServerType]
		if creator == nil {
			logger.Fatalf("创建者为空，配置:%+v", cfg)
			return
		}
		server, err := creator(cfg)
		if err != nil {
			logger.Fatalf("创建服务失败，配置:%+v", cfg)
			return
		}
		server.Init()
		launched[cfg.ServerId] = server
	}

	for _, cfg := range launchArr {
		if server, ok := launched[cfg.ServerId]; ok {
			server.AfterInit()
		}
	}
}

//Shutdown 关闭服务器
func Shutdown() {
	//stop creators
	sort.Slice(launchArr, func(i, j int) bool {
		return launchArr[i].ShutWeight < launchArr[j].ShutWeight
	})
	for _, cfg := range launchArr {
		if server, ok := launched[cfg.ServerId]; ok {
			server.BeforeShutdown()
		}
	}
	for _, cfg := range launchArr {
		if server, ok := launched[cfg.ServerId]; ok {
			server.Shutdown()
		}
	}
	//stop logger writer
	logger.Shutdown()
}
