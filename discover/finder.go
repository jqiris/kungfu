package discover

import (
	"fmt"
	"sync"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/treaty"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Findkey interface {
	int64 | string
}

type ServerMap[k Findkey] map[k]*treaty.Server

type Finder struct {
	serversInt map[string]ServerMap[int64]
	serversStr map[string]ServerMap[string]
	serverLock *sync.RWMutex
}

func NewFinder() *Finder {
	f := &Finder{
		serversInt: make(map[string]ServerMap[int64]),
		serversStr: make(map[string]ServerMap[string]),
		serverLock: new(sync.RWMutex),
	}
	RegServerEventHandlers(f.ServerEventHandler)
	return f
}

func (f *Finder) ServerEventHandler(ev *clientv3.Event, server *treaty.Server) {
	//logger.Infof("server event ev:%+v, server:%+v", ev, server)
	switch ev.Type {
	case clientv3.EventTypePut:
		fallthrough
	case clientv3.EventTypeDelete:
		f.serverLock.Lock()
		delete(f.serversInt, server.ServerType)
		delete(f.serversStr, server.ServerType)
		f.serverLock.Unlock()
	}
}
func (f *Finder) GetServerCache(serverType string, arg any) *treaty.Server {
	f.serverLock.RLock()
	defer f.serverLock.RUnlock()
	switch v := arg.(type) {
	case int64:
		if serverTypeList, ok := f.serversInt[serverType]; ok {
			if server, okv := serverTypeList[v]; okv {
				return server
			}
		}
	case string:
		if serverTypeList, ok := f.serversStr[serverType]; ok {
			if server, okv := serverTypeList[v]; okv {
				return server
			}
		}
	}

	return nil
}

func (f *Finder) GetServerDiscover(serverType string, arg any, options ...FilterOption) *treaty.Server {
	f.serverLock.Lock()
	defer f.serverLock.Unlock()
	server := GetServerByType(serverType, fmt.Sprintf("%v", arg), options...)
	if server != nil {
		switch v := arg.(type) {
		case int64:
			if _, ok := f.serversInt[serverType]; !ok {
				f.serversInt[serverType] = make(ServerMap[int64])
			}
			f.serversInt[serverType][v] = server
			logger.Infof("user server cache,  arg: %v, server_type: %v,server_id:%v", arg, serverType, server.ServerId)
			return server
		case string:
			if _, ok := f.serversStr[serverType]; !ok {
				f.serversStr[serverType] = make(ServerMap[string])
			}
			f.serversStr[serverType][v] = server
			logger.Infof("user server cache,  arg: %v, server_type: %v,server_id:%v", arg, serverType, server.ServerId)
			return server
		}

	}
	return nil
}

func (f *Finder) GetUserServer(serverType string, arg any, options ...FilterOption) *treaty.Server {
	if server := f.GetServerCache(serverType, arg); server != nil {
		return server
	}
	//discover发现
	if server := f.GetServerDiscover(serverType, arg, options...); server != nil {
		return server
	}
	//不存在
	logger.Errorf("找不到服务器：%v", serverType)
	return &treaty.Server{ServerType: "none"}
}

func (f *Finder) RemoveUserCache(arg any) {
	f.serverLock.Lock()
	defer f.serverLock.Unlock()
	switch v := arg.(type) {
	case int64:
		for typ, val := range f.serversInt {
			if _, ok := val[v]; ok {
				delete(f.serversInt[typ], v)
			}
		}
	case string:
		for typ, val := range f.serversStr {
			if _, ok := val[v]; ok {
				delete(f.serversStr[typ], v)
			}
		}
	}
}
