package discover

import (
	"time"

	"stathat.com/c/consistent"

	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/treaty"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.WithField("package", "discover")

	defDiscoverer Discoverer
)

const (
	DiscorverPrefix = "/server/"
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
	Register(server *treaty.Server) error                        //注册服务器
	UnRegister(server *treaty.Server) error                      //注册服务器
	GetServerList() map[string]*treaty.Server                    //获取所有服务信息
	GetServerById(serverId string) *treaty.Server                //根据serverId获取server信息
	GetServerByType(serverType, serverArg string) *treaty.Server //根据serverType及参数分配唯一server信息
}

func Register(server *treaty.Server) error {
	return defDiscoverer.Register(server)
}

func UnRegister(server *treaty.Server) error {
	return defDiscoverer.UnRegister(server)
}

func GetServerList() map[string]*treaty.Server {
	return defDiscoverer.GetServerList()
}

func GetServerById(serverId string) *treaty.Server {
	return defDiscoverer.GetServerById(serverId)
}

func GetServerByType(serverType, serverArg string) *treaty.Server {
	return defDiscoverer.GetServerByType(serverType, serverArg)
}

//serverType stores

type ServerTypeItem struct {
	hash *consistent.Consistent
	list map[string]*treaty.Server
}

func NewServerTypeItem() *ServerTypeItem {
	return &ServerTypeItem{
		hash: consistent.New(),
		list: make(map[string]*treaty.Server),
	}
}
