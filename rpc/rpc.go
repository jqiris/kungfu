package rpc

import (
	"time"

	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/treaty"
	"github.com/nats-io/nats.go"
)

type HttpHandler interface {
	Run(addr ...string) error
}

type CallbackFunc func(req *MsgRpc) []byte

// ServerRpc rpc interface
type ServerRpc interface {
	RegEncoder(typ string, encoder EncoderRpc)                   //register encoder
	Subscribe(s RssBuilder) error                                //self Subscribe
	SubscribeBroadcast(s RssBuilder) error                       //broadcast subscribe
	QueueSubscribe(s RssBuilder) error                           //queue self Subscribe
	Publish(s ReqBuilder) error                                  //publish
	PublishBroadcast(s ReqBuilder) error                         //broadcast publish
	Request(s ReqBuilder) error                                  //request
	Response(codeType string, v interface{}) []byte              //response the msg
	DecodeMsg(codeType string, data []byte, v interface{}) error //decode msg
	GetCoder(codeType string) EncoderRpc                         //get encoder
	GetServer() *treaty.Server                                   //get current server
	Find(serverType string, arg any) *treaty.Server              //find server
	RemoveFindCache(arg any)                                     //clear find cache
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
