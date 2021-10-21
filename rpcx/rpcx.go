package rpcx

import (
	"time"

	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/treaty"
	"github.com/nats-io/nats.go"
)

type CallbackFunc func(req *RpcMsg) []byte

// RpcServer rpc interface
type RpcServer interface {
	Subscribe(server *treaty.Server, callback CallbackFunc, args ...bool) error                        //self Subscribe
	SubscribeJson(server *treaty.Server, callback CallbackFunc, args ...bool) error                    //self json Subscribe
	QueueSubscribe(queue string, server *treaty.Server, callback CallbackFunc, args ...bool) error     //queue self Subscribe
	QueueSubscribeJson(queue string, server *treaty.Server, callback CallbackFunc, args ...bool) error //queue self Subscribe
	Publish(server *treaty.Server, msgId int32, req interface{}) error                                 //publish
	PublishJson(server *treaty.Server, msgId int32, req interface{}) error                             //publish json
	Request(server *treaty.Server, msgId int32, req, resp interface{}) error                           //request
	RequestJson(server *treaty.Server, msgId int32, req, resp interface{}) error                       //request json
	SubscribeBalancer(callback CallbackFunc) error                                                     //balancer subscribe
	SubscribeConnector(callback CallbackFunc) error                                                    //connect subscribe
	SubscribeServer(callback CallbackFunc) error                                                       //server subscribe
	SubscribeDatabase(callback CallbackFunc) error                                                     //database subscribe
	PublishBalancer(msgId int32, req interface{}) error                                                //balancer publish
	PublishConnector(msgId int32, req interface{}) error                                               //connect publish
	PublishServer(msgId int32, req interface{}) error                                                  //server publish
	PublishDatabase(msgId int32, req interface{}) error                                                //database publish
	GetCoder() *RpcEncoder                                                                             //get encoder
	GetJsonCoder() *RpcEncoder                                                                         //get json encoder
	Response(v interface{}) []byte                                                                     //response the msg
	ResponseJson(v interface{}) []byte                                                                 //respon the json msg
	DecodeMsg(data []byte, v interface{}) error                                                        //decode msg
	DecodeJsonMsg(data []byte, v interface{}) error                                                    //decode the  json msg
	GetServer() *treaty.Server                                                                         //get current server
	Find(serverType string, userId int) *treaty.Server                                                 //find server
	Find2(serverType string, arg string) *treaty.Server                                                //find server2
	RemoveFindCache(userId int)                                                                        //clear find cache
	RemoveFindCache2(arg string)                                                                       //clear find cache2
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
