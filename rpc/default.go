package rpc

import (
	"sync"

	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/discover"
	"github.com/jqiris/kungfu/v2/treaty"
)

var (
	defRpc    ServerRpc
	onceInit  sync.Once
	onceClose sync.Once
)

func defRpcInit() {
	onceInit.Do(func() {
		defRpc = NewRpcServer(config.GetRpcConf(), nil)
	})
}

func DefRpcClose() {
	onceClose.Do(func() {
		defRpc.Close()
	})
}

// 公用调用方法
func DefRpcInit() {
	defRpcInit()
}

func Publish(s ReqBuilder) error {
	return defRpc.Publish(s)
}
func QueuePublish(s ReqBuilder) error {
	return defRpc.QueuePublish(s)
}
func PublishBroadcast(s ReqBuilder) error {
	return defRpc.PublishBroadcast(s)
}
func Request(s ReqBuilder) error {
	return defRpc.Request(s)
}
func QueueRequest(s ReqBuilder) error {
	return defRpc.QueueRequest(s)
}

func Find(serverType string, arg any, options ...discover.FilterOption) *treaty.Server {
	return defRpc.Find(serverType, arg, options...)
}

func RemoveFindCache(arg any) {
	defRpc.RemoveFindCache(arg)
}
