package discover

import (
	"context"
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
	if resp, err := kv.Put(context.TODO(), "/server/"+server.RegId(), server.String()); err != nil {
		return err
	} else {
		logger.Infof("EtcdDiscoverer register resp:%+v", resp)
	}
	return nil
}

func (e *EtcdDiscoverer) UnRegister(server common.Server) error {
	kv := clientv3.NewKV(e.Client)
	if resp, err := kv.Delete(context.TODO(), "/server/"+server.RegId(), clientv3.WithPrevKV()); err != nil {
		return err
	} else {
		logger.Infof("EtcdDiscoverer unregister resp:%+v", resp)
	}
	return nil
}

func (e *EtcdDiscoverer) DiscoverServer(serverType common.ServerType) []common.Server {
	kv := clientv3.NewKV(e.Client)
	if resp, err := kv.Get(context.TODO(), fmt.Sprintf("/server/%d/", serverType), clientv3.WithPrefix()); err != nil {
		logger.Errorf("EtcdDiscoverer DiscoverServer err:%v", err)
		return nil
	} else {
		logger.Errorf("EtcdDiscoverer DiscoverServer resp:%v", resp)
	}
	return nil
}

func (e *EtcdDiscoverer) DiscoverServerList() []common.Server {
	kv := clientv3.NewKV(e.Client)
	if resp, err := kv.Get(context.TODO(), "/server/", clientv3.WithPrefix()); err != nil {
		logger.Errorf("EtcdDiscoverer DiscoverServerList err:%v", err)
		return nil
	} else {
		logger.Errorf("EtcdDiscoverer DiscoverServerList resp:%v", resp)
	}
	return nil
}
