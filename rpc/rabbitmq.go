package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/jqiris/kungfu/v2/discover"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/serialize"
	"github.com/jqiris/kungfu/v2/treaty"
	"github.com/jqiris/kungfu/v2/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqRpc struct {
	Endpoints   []string //地址取第一条
	DebugMsg    bool
	Prefix      string
	RpcCoder    map[string]EncoderRpc
	Server      *treaty.Server
	Finder      *discover.Finder
	Client      *amqp.Connection
	DialTimeout time.Duration
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
func WithRabbitMqDialTimeout(timeout time.Duration) RabbitMqRpcOption {
	return func(r *RabbitMqRpc) {
		r.DialTimeout = timeout
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

func NewRpcRabbitMq(opts ...RabbitMqRpcOption) *RabbitMqRpc {
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
	// r.Finder = discover.NewFinder()
	return r
}

func (r *RabbitMqRpc) RegEncoder(typ string, encoder EncoderRpc) {
	if _, ok := r.RpcCoder[typ]; !ok {
		r.RpcCoder[typ] = encoder
	} else {
		logger.Fatalf("encoder type has exist:%v", typ)
	}
}

func (r *RabbitMqRpc) DealMsg(msg amqp.Delivery, callback CallbackFunc, coder EncoderRpc) {
	req := &MsgRpc{}
	err := coder.Decode(msg.Body, req)
	if err != nil {
		logger.Error(err)
		return
	}
	resp := callback(req)
	if resp != nil {
		logger.Info(resp)
	}
	if r.DebugMsg {
		logger.Infof("DealMsg,msgType: %v, msgId: %v", req.MsgType, req.MsgId)
	}
}
func (r *RabbitMqRpc) Subscribe(s RssBuilder) error {
	ch, err := r.Client.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	err = r.prepareMq(ch, s.exName, s.exType, s.queue, s.rtKey)
	if err != nil {
		return err
	}
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	msgs, err := ch.Consume(s.queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go utils.SafeRun(func() {
		for msg := range msgs {
			r.DealMsg(msg, s.callback, coder)
		}
	})
	return nil
}

func (r *RabbitMqRpc) SubscribeBroadcast(s RssBuilder) error {
	panic("not implemented") // TODO: Implement
}

func (r *RabbitMqRpc) QueueSubscribe(s RssBuilder) error {
	panic("not implemented") // TODO: Implement
}

// 准备mq
func (r *RabbitMqRpc) prepareMq(ch *amqp.Channel, exName, exType, queue, rtKey string) error {
	if len(exName) > 0 {
		err := ch.ExchangeDeclare(exName, exType, true, false, false, false, nil)
		if err != nil {
			return err
		}
	}
	_, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	// 绑定任务
	if len(exName) > 0 && len(rtKey) > 0 {
		err := ch.QueueBind(queue, rtKey, exName, false, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RabbitMqRpc) EncodeMsg(coder EncoderRpc, msgType MessageType, msgId int32, req any) ([]byte, error) {
	rpcMsg := &MsgRpc{
		MsgType: msgType,
		MsgId:   msgId,
		MsgData: req,
	}
	data, err := coder.Encode(rpcMsg)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RabbitMqRpc) dialTimeout(s ReqBuilder) time.Duration {
	if s.dialTimeout > 0 {
		return s.dialTimeout
	}
	return r.DialTimeout
}

// 发送消息
func (r *RabbitMqRpc) Publish(s ReqBuilder) error {
	ch, err := r.Client.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	err = r.prepareMq(ch, s.exName, s.exType, s.queue, s.rtKey)
	if err != nil {
		return err
	}
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	data, err := r.EncodeMsg(coder, MsgTypePublish, s.msgId, s.req)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), r.dialTimeout(s))
	defer cancel()
	if len(s.exName) > 0 && len(s.rtKey) > 0 {
		err = ch.PublishWithContext(
			ctx,
			s.exName,
			s.rtKey,
			false,
			false,
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         data,
				DeliveryMode: 2,
			})
	} else {
		err = ch.PublishWithContext(
			ctx,
			DefaultExName,
			s.queue,
			false,
			false,
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         data,
				DeliveryMode: 2,
			})
	}
	return err
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
	return r.RpcCoder[codeType]
}

func (r *RabbitMqRpc) GetServer() *treaty.Server {
	return r.Server
}

func (r *RabbitMqRpc) Find(serverType string, arg any, options ...discover.FilterOption) *treaty.Server {
	// return r.Finder.GetUserServer(serverType, arg, options...)
	return nil
}

func (r *RabbitMqRpc) RemoveFindCache(arg any) {
	// r.Finder.RemoveUserCache(arg)
}

func (r *RabbitMqRpc) Close() error {
	return r.Client.Close()
}
