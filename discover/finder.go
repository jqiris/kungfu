package discover

import (
	"fmt"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/treaty"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
)

type FindKey interface {
	int | string
}

type ServerMap[K FindKey] map[K]*treaty.Server

type Finder struct {
	servers    map[string]ServerMap
	servers2   map[string]map[string]*treaty.Server
	serverLock *sync.RWMutex
}

func NewFinder() *Finder {
	f := &Finder{
		servers:    make(map[string]ServerMap),
		serverLock: new(sync.RWMutex),
	}
	RegEventHandlers(f.ServerEventHandler)
	return f
}

func (f *Finder) ServerEventHandler(ev *clientv3.Event, server *treaty.Server) {
	logger.Infof("server event ev:%+v, server:%+v", ev, server)
	switch ev.Type {
	case clientv3.EventTypePut:
		fallthrough
	case clientv3.EventTypeDelete:
		f.serverLock.Lock()
		delete(f.servers, server.ServerType)
		f.serverLock.Unlock()
	}
}
func (f *Finder) GetServerCache(serverType string, arg any) *treaty.Server {
	f.serverLock.RLock()
	defer f.serverLock.RUnlock()
	if serverTypeList, ok := f.servers[serverType]; ok {
		if server, okv := serverTypeList[arg]; okv {
			return server
		}
	}
	return nil
}

func (f *Finder) GetServerDiscover(serverType string, arg any) *treaty.Server {
	f.serverLock.Lock()
	defer f.serverLock.Unlock()
	server := GetServerByType(serverType, fmt.Sprintf("%v", arg))
	if server != nil {
		if _, ok := f.servers[serverType]; !ok {
			f.servers[serverType] = make(ServerMap)
		}
		f.servers[serverType][arg] = server
		logger.Infof("user server cache,  arg: %v, server_type: %v,server_id:%v", arg, serverType, server.ServerId)
		return server
	}
	return nil
}

func (f *Finder) GetUserServer(serverType string, arg any) *treaty.Server {
	if server := f.GetServerCache(serverType, arg); server != nil {
		return server
	}
	//discover发现
	if server := f.GetServerDiscover(serverType, arg); server != nil {
		return server
	}
	//不存在
	logger.Errorf("找不到服务器：%v", serverType)
	return &treaty.Server{ServerType: "none"}
}

func (f *Finder) RemoveUserCache(arg any) {
	f.serverLock.Lock()
	defer f.serverLock.Unlock()
	for typ, v := range f.servers {
		if _, ok := v[arg]; ok {
			delete(f.servers[typ], arg)
		}
	}
}
