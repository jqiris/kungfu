package discover

import (
	"fmt"
	"github.com/jqiris/kungfu/treaty"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
)

type Finder struct {
	servers    map[string]map[int]*treaty.Server
	serverLock *sync.RWMutex
}

func NewFinder() *Finder {
	f := &Finder{
		servers:    make(map[string]map[int]*treaty.Server),
		serverLock: new(sync.RWMutex),
	}
	RegEventHandlers(f.ServerEventHandler)
	return f
}

func (f *Finder) ServerEventHandler(ev *clientv3.Event, server *treaty.Server) {
	logger.Infof("server event ev:%+v, server:%+v", ev, server)
	switch ev.Type {
	case clientv3.EventTypeDelete:
		f.serverLock.Lock()
		delete(f.servers, server.ServerType)
		f.serverLock.Unlock()
	}
}
func (f *Finder) GetServerCache(serverType string, userId int) *treaty.Server {
	f.serverLock.RLock()
	defer f.serverLock.RUnlock()
	if serverTypeList, ok := f.servers[serverType]; ok {
		if server, okv := serverTypeList[userId]; okv {
			return server
		}
	}
	return nil
}

func (f *Finder) GetServerDiscover(serverType string, userId int) *treaty.Server {
	f.serverLock.Lock()
	defer f.serverLock.Unlock()
	server := GetServerByType(serverType, fmt.Sprintf("%d", userId))
	if server != nil {
		if _, ok := f.servers[serverType]; !ok {
			f.servers[serverType] = make(map[int]*treaty.Server)
		}
		f.servers[serverType][userId] = server
		return server
	}
	return nil
}

func (f *Finder) GetUserServer(serverType string, userId int) *treaty.Server {
	if server := f.GetServerCache(serverType, userId); server != nil {
		return server
	}
	//discover发现
	if server := f.GetServerDiscover(serverType, userId); server != nil {
		return server
	}
	//不存在
	logger.Errorf("找不到服务器：%v", serverType)
	return &treaty.Server{ServerType: "none"}
}
