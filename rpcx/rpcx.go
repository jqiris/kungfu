package rpcx

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/treaty"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	logger = logrus.WithField("package", "rpcx")
)

type CallbackFunc func(coder *RpcEncoder, req *RpcMsg) []byte

//rpc interface
type RpcBase interface {
	Subscribe(server *treaty.Server, callback CallbackFunc) error            //self Subscribe
	Publish(server *treaty.Server, msgId int32, req interface{}) error       //publish
	Request(server *treaty.Server, msgId int32, req, resp interface{}) error //request
}

type RpcBalancer interface {
	RpcBase
	SubscribeBalancer(callback CallbackFunc) error       //gate subscribe
	PublishConnector(msgId int32, req interface{}) error //connect publish
	PublishServer(msgId int32, req interface{}) error    //server publish
}

type RpcConnector interface {
	RpcBase
	SubscribeConnector(callback CallbackFunc) error     //connect subscribe
	PublishBalancer(msgId int32, req interface{}) error //gate publish
	PublishServer(msgId int32, req interface{}) error   //server publish
}

type RpcServer interface {
	RpcBase
	SubscribeServer(callback CallbackFunc) error         //server subscribe
	PublishConnector(msgId int32, req interface{}) error //connect publish
	PublishServer(msgId int32, req interface{}) error    //server publish
}

//create rpc gate
func NewRpcBalancer(cfg config.RpcXConf) RpcBalancer {
	timeout := time.Duration(cfg.DialTimeout) * time.Second
	var r RpcBalancer
	switch cfg.UseType {
	case "nats":
		r = NewRpcNats(WithNatsEndpoints(cfg.Endpoints), WithNatsDialTimeout(timeout), WithNatsOptions(nats.Timeout(timeout)))
	default:
		logger.Fatal("NewRpcBalancer failed")
	}
	return r
}

//create rpc gate
func NewRpcConnector(cfg config.RpcXConf) RpcConnector {
	timeout := time.Duration(cfg.DialTimeout) * time.Second
	var r RpcConnector
	switch cfg.UseType {
	case "nats":
		r = NewRpcNats(WithNatsEndpoints(cfg.Endpoints), WithNatsDialTimeout(timeout), WithNatsOptions(nats.Timeout(timeout)))
	default:
		logger.Fatal("NewRpcConnector failed")
	}
	return r
}

//create rpc server
func NewRpcServer(cfg config.RpcXConf) RpcServer {
	timeout := time.Duration(cfg.DialTimeout) * time.Second
	var r RpcServer
	switch cfg.UseType {
	case "nats":
		r = NewRpcNats(WithNatsEndpoints(cfg.Endpoints), WithNatsDialTimeout(timeout), WithNatsOptions(nats.Timeout(timeout)))
	default:
		logger.Fatal("NewRpcConnector failed")
	}
	return r
}
