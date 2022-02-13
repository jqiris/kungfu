package rpc

import (
	"time"

	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/treaty"
	"github.com/nats-io/nats.go"
)

type CallbackFunc func(req *MsgRpc) []byte

// ServerRpc rpc interface
type ServerRpc interface {
	RegEncoder(typ string, encoder EncoderRpc)                                                        //register encoder
	Subscribe(s SubscriberRpc) error                                                                  //self Subscribe
	SubscribeBroadcast(s SubscriberRpc) error                                                         //broadcast subscribe
	QueueSubscribe(s SubscriberRpc) error                                                             //queue self Subscribe
	Publish(codeType, suffix string, server *treaty.Server, msgId int32, req interface{}) error       //publish
	PublishBalancer(codeType, suffix string, msgId int32, req interface{}) error                      //balancer publish
	PublishConnector(codeType, suffix string, msgId int32, req interface{}) error                     //connect publish
	PublishServer(codeType, suffix string, msgId int32, req interface{}) error                        //server publish
	PublishDatabase(codeType, suffix string, msgId int32, req interface{}) error                      //database publish
	Request(codeType, suffix string, server *treaty.Server, msgId int32, req, resp interface{}) error //request
	Response(codeType string, v interface{}) []byte                                                   //response the msg
	DecodeMsg(codeType string, data []byte, v interface{}) error                                      //decode msg
	GetCoder(codeType string) EncoderRpc                                                              //get encoder
	GetServer() *treaty.Server                                                                        //get current server
	Find(serverType string, arg int) *treaty.Server                                                   //find server
	RemoveFindCache(arg int)                                                                          //clear find cache
}

// NewRpcServer create rpc server
func NewRpcServer(cfg config.RpcConf, server *treaty.Server) ServerRpc {
	timeout := time.Duration(cfg.DialTimeout) * time.Second
	var r ServerRpc
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
