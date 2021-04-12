package discover

import (
	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/treaty"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	logger        = logrus.WithField("package", "discover")
	defDiscoverer Discoverer
)

func InitDiscoverer(cfg conf.DiscoverConf) {
	switch cfg.UseType {
	case "etcd":
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
	Register(server *treaty.Server) error                     //注册服务器
	UnRegister(server *treaty.Server) error                   //注册服务器
	FindServer(serverType treaty.ServerType) []*treaty.Server //获取某个类型的服务器信息
	FindServerList() map[treaty.ServerType][]*treaty.Server   //获取所有的服务器信息
}

func Register(server *treaty.Server) error {
	return defDiscoverer.Register(server)
}

func UnRegister(server *treaty.Server) error {
	return defDiscoverer.UnRegister(server)
}

func FindServerList() map[treaty.ServerType][]*treaty.Server {
	return defDiscoverer.FindServerList()
}

func FindServer(serverType treaty.ServerType) []*treaty.Server {
	return defDiscoverer.FindServer(serverType)
}
