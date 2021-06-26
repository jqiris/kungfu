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

type CallbackFunc func(server RpcServer, req *RpcMsg) []byte

//rpc interface
type RpcServer interface {
	Subscribe(server *treaty.Server, callback CallbackFunc) error            //self Subscribe
	Publish(server *treaty.Server, msgId int32, req interface{}) error       //publish
	Request(server *treaty.Server, msgId int32, req, resp interface{}) error //request
	SubscribeBalancer(callback CallbackFunc) error                           //balancer subscribe
	SubscribeConnector(callback CallbackFunc) error                          //connect subscribe
	SubscribeServer(callback CallbackFunc) error                             //server subscribe
	PublishBalancer(msgId int32, req interface{}) error                      //balancer publish
	PublishConnector(msgId int32, req interface{}) error                     //connect publish
	PublishServer(msgId int32, req interface{}) error                        //server publish
	GetCoder() *RpcEncoder                                                   //get encoder
	Response(v interface{}) []byte                                           //response the msg
	DecodeMsg(data []byte, v interface{}) error                              //decode msg
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
