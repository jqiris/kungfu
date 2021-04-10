package discover

import (
	"github.com/jqiris/kungfu/common"
	"github.com/jqiris/kungfu/conf"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	logger        = logrus.WithField("package", "discover")
	defDiscoverer Discoverer
)

func InitDiscoverer(cfg conf.DiscoverConf) {
	switch cfg.UseType {
	case common.DiscoverEtcd:
		defDiscoverer = NewEtcdDiscoverer(
			WithEtcdDialTimeOut(time.Duration(cfg.DialTimeout)*time.Second),
			WithEtcdEndpoints(cfg.Endpoints),
		)
	default:
		logger.Fatal("InitDiscoverer failed")
	}
}

//find service role
type Discoverer interface {
	Register(server common.Server) error                         //注册服务器
	UnRegister(server common.Server) error                       //注册服务器
	DiscoverServer(serverType common.ServerType) []common.Server //获取某个类型的服务器信息
	DiscoverServerList() map[common.ServerType][]common.Server   //获取所有的服务器信息
}

func Register(server common.Server) error {
	return defDiscoverer.Register(server)
}

func UnRegister(server common.Server) error {
	return defDiscoverer.UnRegister(server)
}

func DiscoverServerList() map[common.ServerType][]common.Server {
	return defDiscoverer.DiscoverServerList()
}

func DiscoverServer(serverType common.ServerType) []common.Server {
	return defDiscoverer.DiscoverServer(serverType)
}
