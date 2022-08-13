package rpc

import (
	"github.com/jqiris/kungfu/v2/discover"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/serialize"
	"github.com/jqiris/kungfu/v2/treaty"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqRpc struct {
	Endpoints []string //地址取第一条
	DebugMsg  bool
	Prefix    string
	RpcCoder  map[string]EncoderRpc
	Server    *treaty.Server
	Finder    *discover.Finder
	Client    *amqp.Connection
}

type RabbitMqRpcOption func(r *RabbitMqRpc)

func WithRabbitMqDebugMsg(debug bool) RabbitMqRpcOption {
	return func(r *RabbitMqRpc) {
		r.DebugMsg = debug
	}
}
func WithRabbitMqEndpoints(endpoints []string) RabbitMqRpcOption {
	return func(r *RabbitMqRpc) {
		r.Endpoints = endpoints
	}
}

func WithRabbitMqServer(server *treaty.Server) RabbitMqRpcOption {
	return func(r *RabbitMqRpc) {
		r.Server = server
	}
}

func WithRabbitMqPrefix(prefix string) RabbitMqRpcOption {
	return func(r *RabbitMqRpc) {
		r.Prefix = prefix
	}
}

func NewRabbitMqRpc(opts ...RabbitMqRpcOption) *RabbitMqRpc {
	r := &RabbitMqRpc{
		Prefix: "rmRpc",
	}
	for _, opt := range opts {
		opt(r)
	}
	if len(r.Endpoints) < 1 {
		logger.Fatal("please set rpc endPoints")
	}
	conn, err := amqp.Dial(r.Endpoints[0])
	if err != nil {
		logger.Fatal(err)
	}
	r.Client = conn
	r.RpcCoder = map[string]EncoderRpc{
		CodeTypeProto: NewRpcEncoder(serialize.NewProtoSerializer()),
		CodeTypeJson:  NewRpcEncoder(serialize.NewJsonSerializer()),
	}
	r.Finder = discover.NewFinder()
	return r
}

func (r *RabbitMqRpc) RegEncoder(typ string, encoder EncoderRpc) {
	if _, ok := r.RpcCoder[typ]; !ok {
		r.RpcCoder[typ] = encoder
	} else {
		logger.Fatalf("encoder type has exist:%v", typ)
	}
}

func (r *RabbitMqRpc) Subscribe(s RssBuilder) error {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) SubscribeBroadcast(s RssBuilder) error {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) QueueSubscribe(s RssBuilder) error {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) Publish(s ReqBuilder) error {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) QueuePublish(s ReqBuilder) error {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) PublishBroadcast(s ReqBuilder) error {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) Request(s ReqBuilder) error {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) QueueRequest(s ReqBuilder) error {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) Response(codeType string, v any) []byte {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) DecodeMsg(codeType string, data []byte, v any) error {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) GetCoder(codeType string) EncoderRpc {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) GetServer() *treaty.Server {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) Find(serverType string, arg any, options ...discover.FilterOption) *treaty.Server {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) RemoveFindCache(arg any) {
	panic("not implemented") // TODO: Implement
}
