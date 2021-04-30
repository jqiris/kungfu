package discover

import (
	"context"
	"time"

	"github.com/jqiris/kungfu/treaty"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdDiscoverer etcd discoverer
type EtcdDiscoverer struct {
	Config        clientv3.Config
	Client        *clientv3.Client
	ServerList    map[string]*treaty.Server
	ServerTypeMap map[string]ServerTypeItem
}
type EtcdOption func(e *EtcdDiscoverer)

func WithEtcdEndpoints(endpoints []string) EtcdOption {
	return func(e *EtcdDiscoverer) {
		e.Config.Endpoints = endpoints
	}
}

func WithEtcdDialTimeOut(d time.Duration) EtcdOption {
	return func(e *EtcdDiscoverer) {
		e.Config.DialTimeout = d
	}
}

// NewEtcdDiscoverer init EtcdDiscoverer
func NewEtcdDiscoverer(opts ...EtcdOption) *EtcdDiscoverer {
	e := &EtcdDiscoverer{
		ServerList:    make(map[string]*treaty.Server),
		ServerTypeMap: make(map[string]ServerTypeItem),
	}
	for _, opt := range opts {
		opt(e)
	}
	cli, err := clientv3.New(e.Config)
	if err != nil {
		logger.Fatal(err)
		return nil
	}
	e.Client = cli
	return e
}

// Init init
func (e *EtcdDiscoverer) Init() {
	go e.Watcher()
	list := e.FindServerList()
	if len(list) > 0 {

	}
}

func (e *EtcdDiscoverer) Watcher() {
	for {
		rch := e.Client.Watch(context.Background(), DiscorverPrefix, clientv3.WithPrefix())
		var err error
		for wresp := range rch {
			err = wresp.Err()
			if err != nil {
				logger.Errorf("etcd watch err:%v", err)
			}
			for _, ev := range wresp.Events {
				logger.Infof("%s %q %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
}

// Register register
func (e *EtcdDiscoverer) Register(server *treaty.Server) error {
	kv := clientv3.NewKV(e.Client)
	ctx, cancel := context.WithTimeout(context.TODO(), e.Config.DialTimeout)
	defer cancel()
	key, val := DiscorverPrefix+treaty.RegSeverItem(server), treaty.RegSerialize(server)
	logger.Infof("discover Register server,k=>v,%s=>%s", key, val)
	if resp, err := kv.Put(ctx, key, val); err != nil {
		return err
	} else {
		logger.Infof("EtcdDiscoverer register resp:%+v", resp)
	}
	return nil
}

func (e *EtcdDiscoverer) UnRegister(server *treaty.Server) error {
	kv := clientv3.NewKV(e.Client)
	ctx, cancel := context.WithTimeout(context.TODO(), e.Config.DialTimeout)
	defer cancel()
	if resp, err := kv.Delete(ctx, DiscorverPrefix+treaty.RegSeverItem(server), clientv3.WithPrevKV()); err != nil {
		return err
	} else {
		logger.Infof("EtcdDiscoverer unregister resp:%+v", resp)
	}
	return nil
}

func (e *EtcdDiscoverer) FindServer(serverType string) []*treaty.Server {
	kv := clientv3.NewKV(e.Client)
	ctx, cancel := context.WithTimeout(context.TODO(), e.Config.DialTimeout)
	defer cancel()
	if resp, err := kv.Get(ctx, DiscorverPrefix+serverType+"/", clientv3.WithPrefix()); err != nil {
		logger.Errorf("EtcdDiscoverer FindServer err:%v", err)
		return nil
	} else {
		if resp.Count > 0 {
			res := make([]*treaty.Server, 0)
			for _, v := range resp.Kvs {
				if server, err := treaty.RegUnSerialize(v.Value); err == nil {
					res = append(res, server)
				} else {
					logger.Errorf("EtcdDiscoverer FindServer err:%+v", err)
				}
			}
			return res
		}
	}
	return nil
}

func (e *EtcdDiscoverer) FindServerList() map[string][]*treaty.Server {
	kv := clientv3.NewKV(e.Client)
	ctx, cancel := context.WithTimeout(context.TODO(), e.Config.DialTimeout)
	defer cancel()
	if resp, err := kv.Get(ctx, DiscorverPrefix, clientv3.WithPrefix()); err != nil {
		logger.Errorf("EtcdDiscoverer FindServerList err:%v", err)
		return nil
	} else {
		if resp.Count > 0 {
			res := make(map[string][]*treaty.Server)
			for _, v := range resp.Kvs {
				if server, err := treaty.RegUnSerialize(v.Value); err == nil {
					res[server.ServerType] = append(res[server.ServerType], server)
				} else {
					logger.Errorf("EtcdDiscoverer FindServerList err:%+v", err)
				}
			}
			return res
		}
	}
	return nil
}
