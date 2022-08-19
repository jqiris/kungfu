package rpc

import (
	"context"
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jqiris/kungfu/v2/discover"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/serialize"
	"github.com/jqiris/kungfu/v2/treaty"
	"github.com/jqiris/kungfu/v2/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitWaitItem struct {
	CorrId   string
	CodeType string
	MsgData  any
	MsgReply chan *MsgRpc
}

type RabbitReplyQueue struct {
	QueueName string
	WaitMap   map[string]*RabbitWaitItem
	WaitChan  chan *RabbitWaitItem
	ReplyChan <-chan amqp.Delivery
	RpcCoder  map[string]EncoderRpc
}

func NewRabbitReplyQueue(name string, ch <-chan amqp.Delivery, rpcCoder map[string]EncoderRpc) *RabbitReplyQueue {
	return &RabbitReplyQueue{
		QueueName: name,
		WaitMap:   make(map[string]*RabbitWaitItem),
		WaitChan:  make(chan *RabbitWaitItem, 30),
		ReplyChan: ch,
		RpcCoder:  rpcCoder,
	}
}

func (r *RabbitReplyQueue) WaitReply() {
	go utils.SafeRun(func() {
		for {
			select {
			case item := <-r.WaitChan:
				r.WaitMap[item.CorrId] = item
			case reply := <-r.ReplyChan:
				if v, ok := r.WaitMap[reply.CorrelationId]; ok {
					coder := r.RpcCoder[v.CodeType]
					if coder == nil {
						logger.Errorf("rpc coder not exist:%v", v.CodeType)
						continue
					}
					respMsg := &MsgRpc{MsgData: v.MsgData}
					err := coder.Decode(reply.Body, respMsg)
					if err != nil {
						logger.Error(err)
					}
					v.MsgReply <- respMsg
					delete(r.WaitMap, reply.CorrelationId)
				} else {
					logger.Errorf("WaitReply can't find reply msg,queue:%v,corrid:%v", r.QueueName, reply.CorrelationId)
				}
			}
		}
	})
}

func (r *RabbitMqRpc) GetReplyQueue(subReply string) (*RabbitReplyQueue, error) {
	if v, ok := r.ReplyQueues.Load(subReply); ok {
		return v.(*RabbitReplyQueue), nil
	}
	conn, err := r.OpenConn()
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	replyCh, err := ch.QueueDeclare(
		subReply, // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // noWait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}
	err = ch.Qos(
		30,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}
	msgs, err := ch.Consume(
		replyCh.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return nil, err
	}
	queue := NewRabbitReplyQueue(subReply, msgs, r.RpcCoder)
	queue.WaitReply()
	r.ReplyQueues.Store(subReply, queue)
	return queue, nil
}

type RabbitMqRpc struct {
	Endpoints   []string //地址取第一条
	DebugMsg    bool
	Prefix      string
	RpcCoder    map[string]EncoderRpc
	Server      *treaty.Server
	Finder      *discover.Finder
	Client      *amqp.Connection
	DialTimeout time.Duration
	ReplyQueues sync.Map
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
		Prefix:      "rmRpc",
		ReplyQueues: sync.Map{},
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

func (r *RabbitMqRpc) OpenConn() (*amqp.Connection, error) {
	return amqp.Dial(r.Endpoints[0])
}

func (r *RabbitMqRpc) RegEncoder(typ string, encoder EncoderRpc) {
	if _, ok := r.RpcCoder[typ]; !ok {
		r.RpcCoder[typ] = encoder
	} else {
		logger.Fatalf("encoder type has exist:%v", typ)
	}
}

func (r *RabbitMqRpc) DealMsg(s RssBuilder, ch *amqp.Channel, msg amqp.Delivery, callback CallbackFunc, coder EncoderRpc) {
	req := &MsgRpc{}
	err := coder.Decode(msg.Body, req)
	if err != nil {
		logger.Error(err)
		return
	}
	resp := callback(req)
	if resp != nil {
		replyCh, err := r.Client.Channel()
		if err != nil {
			logger.Errorf("DealMsg 回复报错:subReply:%v,corrid:%v,err:%v", msg.ReplyTo, msg.CorrelationId, err)
			return
		}
		defer replyCh.Close()
		dialTimeout := r.dialTimeoutRss(s)
		ctx, cancel := context.WithTimeout(context.TODO(), dialTimeout)
		defer cancel()
		if err = replyCh.PublishWithContext(
			ctx,
			"",          // exchange
			msg.ReplyTo, // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: msg.CorrelationId,
				Body:          resp,
				DeliveryMode:  2,
			}); err != nil {
			logger.Error(err)
		} else if r.DebugMsg {
			logger.Infof("DealMsg 回复消息:subReply:%v,corrid:%v", msg.ReplyTo, msg.CorrelationId)
		}
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
	sub := path.Join(r.Prefix, treaty.RegSeverItem(s.server), s.suffix)
	err = r.prepareMq(ch, s.exName, s.exType, sub, sub)
	if err != nil {
		return err
	}
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	msgs, err := ch.Consume(sub, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go utils.SafeRun(func() {
		for msg := range msgs {
			utils.SafeRun(func() {
				r.DealMsg(s, ch, msg, s.callback, coder)
			})
		}
	})
	return nil
}

func (r *RabbitMqRpc) SubscribeBroadcast(s RssBuilder) error {
	ch, err := r.Client.Channel()
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, s.server.ServerType, s.suffix)
	err = r.prepareMq(ch, s.exName, s.exType, sub, sub)
	if err != nil {
		return err
	}
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	msgs, err := ch.Consume(sub, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go utils.SafeRun(func() {
		for msg := range msgs {
			utils.SafeRun(func() {
				r.DealMsg(s, ch, msg, s.callback, coder)
			})
		}
	})
	return nil
}

func (r *RabbitMqRpc) QueueSubscribe(s RssBuilder) error {
	ch, err := r.Client.Channel()
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, treaty.RegSeverQueue(s.server.ServerType, s.queue), s.suffix)
	err = r.prepareMq(ch, s.exName, s.exType, sub, sub)
	if err != nil {
		return err
	}
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	msgs, err := ch.Consume(sub, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go utils.SafeRun(func() {
		for msg := range msgs {
			utils.SafeRun(func() {
				r.DealMsg(s, ch, msg, s.callback, coder)
			})
		}
	})
	return nil
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

func (r *RabbitMqRpc) EncodeMsgRaw(coder EncoderRpc, msgType MessageType, msgId int32, req any) ([]byte, error) {
	rpcMsg := &MsgRpc{
		MsgType: msgType,
		MsgId:   msgId,
		MsgData: req,
	}
	data, err := coder.EncodeMsg(rpcMsg)
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

func (r *RabbitMqRpc) dialTimeoutRss(s RssBuilder) time.Duration {
	if s.dialTimeout > 0 {
		return s.dialTimeout
	}
	return r.DialTimeout
}

// 发送消息
func (r *RabbitMqRpc) SendMsg(s ReqBuilder) error {
	ch, err := r.Client.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	queue := s.queue
	rtKey := s.rtKey
	if len(r.Prefix) > 0 {
		queue = r.Prefix + "_" + s.queue
	}
	if len(r.Prefix) > 0 && len(s.rtKey) > 0 {
		rtKey = r.Prefix + "_" + s.rtKey
	}
	err = r.prepareMq(ch, s.exName, s.exType, queue, rtKey)
	if err != nil {
		return err
	}
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	data, err := r.EncodeMsgRaw(coder, MsgTypePublish, s.msgId, s.req)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), r.dialTimeout(s))
	defer cancel()
	if len(s.exName) > 0 && len(rtKey) > 0 {
		err = ch.PublishWithContext(
			ctx,
			s.exName,
			rtKey,
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
			"",
			queue,
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

// 发送消息
func (r *RabbitMqRpc) Publish(s ReqBuilder) error {
	ch, err := r.Client.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	sub := path.Join(r.Prefix, treaty.RegSeverItem(s.server), s.suffix)
	err = r.prepareMq(ch, s.exName, s.exType, sub, sub)
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
	if len(s.exName) > 0 && len(sub) > 0 {
		err = ch.PublishWithContext(
			ctx,
			s.exName,
			sub,
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
			sub,
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
	ch, err := r.Client.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	sub := path.Join(r.Prefix, treaty.RegSeverQueue(s.serverType, s.queue), s.suffix)
	err = r.prepareMq(ch, s.exName, s.exType, sub, sub)
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
	if len(s.exName) > 0 && len(sub) > 0 {
		err = ch.PublishWithContext(
			ctx,
			s.exName,
			sub,
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
			sub,
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

func (r *RabbitMqRpc) PublishBroadcast(s ReqBuilder) error {
	ch, err := r.Client.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	sub := path.Join(r.Prefix, s.serverType, s.suffix)
	err = r.prepareMq(ch, s.exName, s.exType, sub, sub)
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
	if len(s.exName) > 0 && len(sub) > 0 {
		err = ch.PublishWithContext(
			ctx,
			s.exName,
			sub,
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
			sub,
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

func (r *RabbitMqRpc) Request(s ReqBuilder) error {
	ch, err := r.Client.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	sub := path.Join(r.Prefix, treaty.RegSeverItem(s.server), s.suffix)
	err = r.prepareMq(ch, s.exName, s.exType, sub, sub)
	if err != nil {
		return err
	}
	corrId := uuid.NewString()
	subReply := path.Join(sub, DefaultReply)
	replyQueue, err := r.GetReplyQueue(subReply)
	if err != nil {
		return err
	}
	replyItem := &RabbitWaitItem{
		CorrId:   corrId,
		CodeType: s.codeType,
		MsgData:  s.resp,
		MsgReply: make(chan *MsgRpc, 1),
	}
	replyQueue.WaitChan <- replyItem
	if r.DebugMsg {
		logger.Infof("Request 发送消息:subReply:%v,corrid:%v", subReply, corrId)
	}
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	data, err := r.EncodeMsg(coder, MsgTypeRequest, s.msgId, s.req)
	if err != nil {
		return err
	}
	dialTimeout := r.dialTimeout(s)
	ctx, cancel := context.WithTimeout(context.TODO(), dialTimeout)
	defer cancel()
	err = ch.PublishWithContext(
		ctx,
		DefaultExName,
		sub,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       subReply,
			Body:          data,
			DeliveryMode:  2,
		})
	if err != nil {
		return err
	}
	replyCtx, replyCancel := context.WithTimeout(context.TODO(), 2*dialTimeout)
	defer replyCancel()
	for {
		select {
		case item := <-replyItem.MsgReply:
			s.resp = item
			if r.DebugMsg {
				logger.Infof("Request 收到消息:subReply:%v,corrid:%v", subReply, corrId)
			}
			return nil
		case <-replyCtx.Done():
			return fmt.Errorf("消息返回超时,subReply:%v,corrId:%v", subReply, corrId)
		}
	}
}

func (r *RabbitMqRpc) QueueRequest(s ReqBuilder) error {
	ch, err := r.Client.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	sub := path.Join(r.Prefix, treaty.RegSeverQueue(s.serverType, s.queue), s.suffix)
	err = r.prepareMq(ch, s.exName, s.exType, sub, sub)
	if err != nil {
		return err
	}
	corrId := uuid.NewString()
	subReply := path.Join(sub, DefaultReply)
	replyQueue, err := r.GetReplyQueue(subReply)
	if err != nil {
		return err
	}
	replyItem := &RabbitWaitItem{
		CorrId:   corrId,
		CodeType: s.codeType,
		MsgData:  s.resp,
		MsgReply: make(chan *MsgRpc, 1),
	}
	replyQueue.WaitChan <- replyItem
	logger.Warnf("QueueRequest 发送消息:subReply:%v,corrid:%v", subReply, corrId)
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	data, err := r.EncodeMsg(coder, MsgTypeRequest, s.msgId, s.req)
	if err != nil {
		return err
	}
	dialTimeout := r.dialTimeout(s)
	ctx, cancel := context.WithTimeout(context.TODO(), dialTimeout)
	defer cancel()
	err = ch.PublishWithContext(
		ctx,
		DefaultExName,
		sub,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       subReply,
			Body:          data,
			DeliveryMode:  2,
		})
	if err != nil {
		return err
	}
	replyCtx, replyCancel := context.WithTimeout(context.TODO(), 2*dialTimeout)
	defer replyCancel()
	for {
		select {
		case item := <-replyItem.MsgReply:
			s.resp = item
			logger.Warnf("QueueRequest 收到消息:subReply:%v,corrid:%v", subReply, corrId)
			return nil
		case <-replyCtx.Done():
			return fmt.Errorf("消息返回超时,subReply:%v,corrId:%v", subReply, corrId)
		}
	}
}

func (r *RabbitMqRpc) Response(codeType string, v any) []byte {
	coder := r.RpcCoder[codeType]
	if coder == nil {
		logger.Errorf("rpc coder not exist:%v", codeType)
		return nil
	}
	return coder.Response(v)
}

func (r *RabbitMqRpc) DecodeMsg(codeType string, data []byte, v any) error {
	coder := r.RpcCoder[codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", codeType)
	}
	return coder.DecodeMsg(data, v)
}

func (r *RabbitMqRpc) GetCoder(codeType string) EncoderRpc {
	return r.RpcCoder[codeType]
}

func (r *RabbitMqRpc) GetServer() *treaty.Server {
	return r.Server
}

func (r *RabbitMqRpc) Find(serverType string, arg any, options ...discover.FilterOption) *treaty.Server {
	return r.Finder.GetUserServer(serverType, arg, options...)
}

func (r *RabbitMqRpc) RemoveFindCache(arg any) {
	r.Finder.RemoveUserCache(arg)
}

func (r *RabbitMqRpc) Close() error {
	return r.Client.Close()
}
