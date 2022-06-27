package discover

import (
	"time"

	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/logger"
	clientv3 "go.etcd.io/etcd/client/v3"

	"stathat.com/c/consistent"

	"github.com/jqiris/kungfu/v2/treaty"
)

var (
	defDiscoverer Discoverer
)

func InitDiscoverer(cfg config.DiscoverConf) {
	switch cfg.UseType {
	case "etcd":
		defDiscoverer = NewEtcdDiscoverer(
			WithEtcdDialTimeOut(time.Duration(cfg.DialTimeout)*time.Second),
			WithEtcdEndpoints(cfg.Endpoints),
			WithEtcdPrefix(cfg.Prefix),
		)
	default:
		logger.Fatal("InitDiscoverer failed")
	}
}

type EventHandler func(ev *clientv3.Event, server *treaty.Server)

//Discoverer find service role
type Discoverer interface {
	Register(server *treaty.Server) error                                      //注册服务器
	UnRegister(server *treaty.Server) error                                    //注册服务器
	GetServerList(args ...bool) map[string]*treaty.Server                      //获取所有服务信息
	GetServerById(serverId string, args ...bool) *treaty.Server                //根据serverId获取server信息
	GetServerByType(serverType, serverArg string, args ...bool) *treaty.Server //根据serverType及参数分配唯一server信息
	GetServerByTypeLoad(serverType string, args ...bool) *treaty.Server        //根绝服务器最小负载量选择服务
	GetServerTypeList(serverType string, args ...bool) map[string]*treaty.Server
	RegEventHandlers(handlers ...EventHandler)
	EventHandlerExec(ev *clientv3.Event, server *treaty.Server)
	IncreLoad(serverId string, load int64) error //负载增加
	DecreLoad(serverId string, load int64) error //负载减少
}

func IncreLoad(serverId string, load int64) error {
	return defDiscoverer.IncreLoad(serverId, load)
}

func DecreLoad(serverId string, load int64) error {
	return defDiscoverer.DecreLoad(serverId, load)
}

func Register(server *treaty.Server) error {
	return defDiscoverer.Register(server)
}

func UnRegister(server *treaty.Server) error {
	return defDiscoverer.UnRegister(server)
}

func GetServerList(args ...bool) map[string]*treaty.Server {
	return defDiscoverer.GetServerList(args...)
}

func GetServerById(serverId string, args ...bool) *treaty.Server {
	return defDiscoverer.GetServerById(serverId, args...)
}

func GetServerByType(serverType, serverArg string, args ...bool) *treaty.Server {
	return defDiscoverer.GetServerByType(serverType, serverArg, args...)
}

func GetServerByTypeLoad(serverType string, args ...bool) *treaty.Server {
	return defDiscoverer.GetServerByTypeLoad(serverType, args...)
}

func GetServerTypeList(serverType string, args ...bool) map[string]*treaty.Server {
	return defDiscoverer.GetServerTypeList(serverType, args...)
}

func RegEventHandlers(handlers ...EventHandler) {
	defDiscoverer.RegEventHandlers(handlers...)
}

//serverType stores

type ServerTypeItem struct {
	hash *consistent.Consistent
	List map[string]*treaty.Server
}

func NewServerTypeItem() *ServerTypeItem {
	return &ServerTypeItem{
		hash: consistent.New(),
		List: make(map[string]*treaty.Server),
	}
}
