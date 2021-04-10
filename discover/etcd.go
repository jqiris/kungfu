package discover

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jqiris/kungfu/common"
	"go.etcd.io/etcd/client/v3"
	"time"
)

//etcd discoverer
type EtcdDiscoverer struct {
	Config clientv3.Config
	Client *clientv3.Client
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

//init EtcdDiscoverer
func NewEtcdDiscoverer(opts ...EtcdOption) *EtcdDiscoverer {
	e := &EtcdDiscoverer{}
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

//register
func (e *EtcdDiscoverer) Register(server common.Server) error {
	if server.ServerId < common.MinServerId {
		return errors.New("ServerId cannot less than MinServerId")
	}
	kv := clientv3.NewKV(e.Client)
	ctx, cancel := context.WithTimeout(context.TODO(), e.Config.DialTimeout)
	defer cancel()
	if resp, err := kv.Put(ctx, "/server/"+server.RegId(), server.String()); err != nil {
		return err
	} else {
		logger.Infof("EtcdDiscoverer register resp:%+v", resp)
	}
	return nil
}

func (e *EtcdDiscoverer) UnRegister(server common.Server) error {
	kv := clientv3.NewKV(e.Client)
	ctx, cancel := context.WithTimeout(context.TODO(), e.Config.DialTimeout)
	defer cancel()
	if resp, err := kv.Delete(ctx, "/server/"+server.RegId(), clientv3.WithPrevKV()); err != nil {
		return err
	} else {
		logger.Infof("EtcdDiscoverer unregister resp:%+v", resp)
	}
	return nil
}

func (e *EtcdDiscoverer) DiscoverServer(serverType common.ServerType) []common.Server {
	kv := clientv3.NewKV(e.Client)
	ctx, cancel := context.WithTimeout(context.TODO(), e.Config.DialTimeout)
	defer cancel()
	if resp, err := kv.Get(ctx, fmt.Sprintf("/server/%d/", serverType), clientv3.WithPrefix()); err != nil {
		logger.Errorf("EtcdDiscoverer DiscoverServer err:%v", err)
		return nil
	} else {
		if resp.Count > 0 {
			res := make([]common.Server, 0)
			for _, v := range resp.Kvs {
				var server common.Server
				if err := json.Unmarshal(v.Value, &server); err == nil {
					res = append(res, server)
				} else {
					logger.Errorf("EtcdDiscoverer DiscoverServer err:%+v", err)
				}
			}
			return res
		}
	}
	return nil
}

func (e *EtcdDiscoverer) DiscoverServerList() map[common.ServerType][]common.Server {
	kv := clientv3.NewKV(e.Client)
	ctx, cancel := context.WithTimeout(context.TODO(), e.Config.DialTimeout)
	defer cancel()
	if resp, err := kv.Get(ctx, "/server/", clientv3.WithPrefix()); err != nil {
		logger.Errorf("EtcdDiscoverer DiscoverServerList err:%v", err)
		return nil
	} else {
		if resp.Count > 0 {
			res := make(map[common.ServerType][]common.Server)
			for _, v := range resp.Kvs {
				var server common.Server
				if err := json.Unmarshal(v.Value, &server); err == nil {
					res[server.ServerType] = append(res[server.ServerType], server)
				} else {
					logger.Errorf("EtcdDiscoverer DiscoverServerList err:%+v", err)
				}
			}
			return res
		}
	}
	return nil
}
