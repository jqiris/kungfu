package launch

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpc"
)

//服务器集群管理
var (
	creators map[string]rpc.ServerCreator
	launched map[string]rpc.ServerEntity
)

func init() {
	creators = make(map[string]rpc.ServerCreator)
	launched = make(map[string]rpc.ServerEntity)
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
	for sid, cfg := range servers {
		if cfg.IsLaunch {
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
			server.AfterInit()
			launched[sid] = server
		}
	}
}

//Shutdown 关闭服务器
func Shutdown() {
	//stop creators
	for _, server := range launched {
		server.BeforeShutdown()
		server.Shutdown()
	}
}
