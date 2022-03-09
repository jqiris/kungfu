package discover

import (
	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"

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
	Register(server *treaty.Server) error                        //注册服务器
	UnRegister(server *treaty.Server) error                      //注册服务器
	GetServerList() map[string]*treaty.Server                    //获取所有服务信息
	GetServerById(serverId string) *treaty.Server                //根据serverId获取server信息
	GetServerByType(serverType, serverArg string) *treaty.Server //根据serverType及参数分配唯一server信息
	GetServerByTypeLoad(serverType string) *treaty.Server        //根绝服务器最小负载量选择服务
	GetServerTypeList(serverType string) map[string]*treaty.Server
	RegEventHandlers(handlers ...EventHandler)
	EventHandlerExec(ev *clientv3.Event, server *treaty.Server)
	IncrLoad(server *treaty.Server) error //增加负载量
	DecrLoad(server *treaty.Server) error //减少负载量
}

func Register(server *treaty.Server) error {
	return defDiscoverer.Register(server)
}

func UnRegister(server *treaty.Server) error {
	return defDiscoverer.UnRegister(server)
}

func IncrLoad(server *treaty.Server) error {
	return defDiscoverer.IncrLoad(server)
}

func DecrLoad(server *treaty.Server) error {
	return defDiscoverer.DecrLoad(server)
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

func GetServerByTypeLoad(serverType string) *treaty.Server {
	return defDiscoverer.GetServerByTypeLoad(serverType)
}

func GetServerTypeList(serverType string) map[string]*treaty.Server {
	return defDiscoverer.GetServerTypeList(serverType)
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
