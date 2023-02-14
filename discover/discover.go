/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

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
			WithEtcdServerPrefix(cfg.ServerPrefix),
			WithEtcdDataPrefix(cfg.DataPrefix),
		)
	default:
		logger.Fatal("InitDiscoverer failed")
	}
}

type ServerEventHandler func(ev *clientv3.Event, server *treaty.Server)
type DataEventHandler func(ev *clientv3.Event)

//Discoverer find service role
type Discoverer interface {
	Register(server *treaty.Server) error                                                 //注册服务器
	UnRegister(server *treaty.Server) error                                               //注册服务器
	GetServerList(options ...FilterOption) map[string]*treaty.Server                      //获取所有服务信息
	GetServerById(serverId string, options ...FilterOption) *treaty.Server                //根据serverId获取server信息
	GetServerByType(serverType, serverArg string, options ...FilterOption) *treaty.Server //根据serverType及参数分配唯一server信息
	GetServerByTypeLoad(serverType string, options ...FilterOption) *treaty.Server        //根绝服务器最小负载量选择服务
	GetServerTypeList(serverType string, options ...FilterOption) map[string]*treaty.Server
	RegServerEventHandlers(handlers ...ServerEventHandler)
	ServerEventHandlerExec(ev *clientv3.Event, server *treaty.Server)
	RegDataEventHandlers(handlers ...DataEventHandler)
	DataEventHandlerExec(ev *clientv3.Event)
	IncreLoad(serverId string, load int64, options ...FilterOption) error //负载增加
	DecreLoad(serverId string, load int64, options ...FilterOption) error //负载减少
	PutData(key, val string) error                                        //增加数据
	RemoveData(key string) error                                          //删除数据
	GetData(key string) (string, error)
}

func PutData(key, val string) error {
	return defDiscoverer.PutData(key, val)
}

func RemoveData(key string) error {
	return defDiscoverer.RemoveData(key)
}

func GetData(key string) (string, error) {
	return defDiscoverer.GetData(key)
}

func IncreLoad(serverId string, load int64, options ...FilterOption) error {
	return defDiscoverer.IncreLoad(serverId, load, options...)
}

func DecreLoad(serverId string, load int64, options ...FilterOption) error {
	return defDiscoverer.DecreLoad(serverId, load, options...)
}

func Register(server *treaty.Server) error {
	return defDiscoverer.Register(server)
}

func UnRegister(server *treaty.Server) error {
	return defDiscoverer.UnRegister(server)
}

func GetServerList(options ...FilterOption) map[string]*treaty.Server {
	return defDiscoverer.GetServerList(options...)
}

func GetServerById(serverId string, options ...FilterOption) *treaty.Server {
	return defDiscoverer.GetServerById(serverId, options...)
}

func GetServerByType(serverType, serverArg string, options ...FilterOption) *treaty.Server {
	return defDiscoverer.GetServerByType(serverType, serverArg, options...)
}

func GetServerByTypeLoad(serverType string, options ...FilterOption) *treaty.Server {
	return defDiscoverer.GetServerByTypeLoad(serverType, options...)
}

func GetServerTypeList(serverType string, options ...FilterOption) map[string]*treaty.Server {
	return defDiscoverer.GetServerTypeList(serverType, options...)
}

func RegServerEventHandlers(handlers ...ServerEventHandler) {
	defDiscoverer.RegServerEventHandlers(handlers...)
}

func RegDataEventHandlers(handlers ...DataEventHandler) {
	defDiscoverer.RegDataEventHandlers(handlers...)
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
