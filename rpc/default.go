package rpc

import (
	"sync"

	"github.com/jqiris/kungfu/v2/config"
)

var (
	defRpc ServerRpc
	once   sync.Once
)

func defRpcInit() {
	once.Do(func() {
		defRpc = NewRpcServer(config.GetRpcConf(), nil)
	})
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
