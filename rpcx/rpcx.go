package rpcx

import (
	"time"

	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/treaty"
	"github.com/nats-io/nats.go"
)

type CallbackFunc func(server RpcServer, req *RpcMsg) []byte

// RpcServer rpc interface
type RpcServer interface {
	Subscribe(server *treaty.Server, callback CallbackFunc) error                    //self Subscribe
	QueueSubscribe(queue string, server *treaty.Server, callback CallbackFunc) error //queue self Subscribe
	Publish(server *treaty.Server, msgId int32, req interface{}) error               //publish
	Request(server *treaty.Server, msgId int32, req, resp interface{}) error         //request
	SubscribeBalancer(callback CallbackFunc) error                                   //balancer subscribe
	SubscribeConnector(callback CallbackFunc) error                                  //connect subscribe
	SubscribeServer(callback CallbackFunc) error                                     //server subscribe
	SubscribeDatabase(callback CallbackFunc) error                                   //database subscribe
	PublishBalancer(msgId int32, req interface{}) error                              //balancer publish
	PublishConnector(msgId int32, req interface{}) error                             //connect publish
	PublishServer(msgId int32, req interface{}) error                                //server publish
	PublishDatabase(msgId int32, req interface{}) error                              //database publish
	GetCoder() *RpcEncoder                                                           //get encoder
	Response(v interface{}) []byte                                                   //response the msg
	DecodeMsg(data []byte, v interface{}) error                                      //decode msg
	GetServer() *treaty.Server                                                       //get current server
}

// NewRpcServer create rpc server
func NewRpcServer(cfg config.RpcXConf, server *treaty.Server) RpcServer {
	timeout := time.Duration(cfg.DialTimeout) * time.Second
	var r RpcServer
	switch cfg.UseType {
	case "nats":
		r = NewRpcNats(
			WithNatsEndpoints(cfg.Endpoints),
			WithNatsDialTimeout(timeout),
			WithNatsOptions(nats.Timeout(timeout)),
			WithNatsServer(server),
			WithNatsPrefix(cfg.Prefix),
			WithNatsDebugMsg(cfg.DebugMsg),
		)
	default:
		logger.Fatal("NewRpcConnector failed")
	}
	return r
}
