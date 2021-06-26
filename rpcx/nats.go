package rpcx

import (
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/treaty"
	"github.com/nats-io/nats.go"
	"path"
	"strings"
	"time"
)

const (
	RpcPrefix = "RpcX"
	Balancer  = "balancer"
	Connector = "connector"
	Server    = "server"
)

type RpcNats struct {
	Endpoints   []string
	Options     []nats.Option
	Client      *nats.Conn
	DialTimeout time.Duration
	RpcCoder    *RpcEncoder
	Server      *treaty.Server
}
type RpcNatsOption func(r *RpcNats)

func WithNatsEndpoints(endpoints []string) RpcNatsOption {
	return func(r *RpcNats) {
		r.Endpoints = endpoints
	}
}
func WithNatsDialTimeout(timeout time.Duration) RpcNatsOption {
	return func(r *RpcNats) {
		r.DialTimeout = timeout
	}
}
func WithNatsServer(server *treaty.Server) RpcNatsOption {
	return func(r *RpcNats) {
		r.Server = server
	}
}
func WithNatsOptions(opts ...nats.Option) RpcNatsOption {
	return func(r *RpcNats) {
		r.Options = opts
	}
}

func NewRpcNats(opts ...RpcNatsOption) *RpcNats {
	r := &RpcNats{}
	for _, opt := range opts {
		opt(r)
	}
	url := strings.Join(r.Endpoints, ",")
	conn, err := nats.Connect(url, r.Options...)
	if err != nil {
		logger.Fatal(err)
	}
	r.Client = conn
	r.RpcCoder = NewRpcEncoder()
	return r
}

func (r *RpcNats) Subscribe(server *treaty.Server, callback CallbackFunc) error {
	sub := path.Join(RpcPrefix, treaty.RegSeverItem(server))
	if _, err := r.Client.Subscribe(sub, func(msg *nats.Msg) {
		go helper.SafeRun(func() {
			r.DealMsg(msg, callback)
		})
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) SubscribeBalancer(callback CallbackFunc) error {
	sub := path.Join(RpcPrefix, Balancer)
	if _, err := r.Client.Subscribe(sub, func(msg *nats.Msg) {
		go helper.SafeRun(func() {
			r.DealMsg(msg, callback)
		})
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) SubscribeConnector(callback CallbackFunc) error {
	sub := path.Join(RpcPrefix, Connector)
	if _, err := r.Client.Subscribe(sub, func(msg *nats.Msg) {
		go helper.SafeRun(func() {
			r.DealMsg(msg, callback)
		})
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) SubscribeServer(callback CallbackFunc) error {
	sub := path.Join(RpcPrefix, Server)
	if _, err := r.Client.Subscribe(sub, func(msg *nats.Msg) {
		go helper.SafeRun(func() {
			r.DealMsg(msg, callback)
		})
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) DealMsg(msg *nats.Msg, callback CallbackFunc) {
	req := &RpcMsg{}
	err := r.RpcCoder.Decode(msg.Data, req)
	if err != nil {
		logger.Error(err)
		return
	}
	resp := callback(r, req)
	if resp != nil {
		if err := msg.Respond(resp); err != nil {
			logger.Error(err)
		}
	}
}

func (r *RpcNats) Request(server *treaty.Server, msgId int32, req, resp interface{}) error {
	var msg *nats.Msg
	var err error
	var data []byte
	data, err = r.EncodeMsg(Request, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(RpcPrefix, treaty.RegSeverItem(server))
	if msg, err = r.Client.Request(sub, data, r.DialTimeout); err == nil {
		respMsg := &RpcMsg{MsgData: resp}
		err = r.RpcCoder.Decode(msg.Data, respMsg)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (r *RpcNats) Publish(server *treaty.Server, msgId int32, req interface{}) error {
	data, err := r.EncodeMsg(Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(RpcPrefix, treaty.RegSeverItem(server))
	if err = r.Client.Publish(sub, data); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) PublishBalancer(msgId int32, req interface{}) error {
	data, err := r.EncodeMsg(Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(RpcPrefix, Balancer)
	return r.Client.Publish(sub, data)
}

func (r *RpcNats) PublishConnector(msgId int32, req interface{}) error {
	data, err := r.EncodeMsg(Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(RpcPrefix, Connector)
	return r.Client.Publish(sub, data)
}

func (r *RpcNats) PublishServer(msgId int32, req interface{}) error {
	data, err := r.EncodeMsg(Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(RpcPrefix, Server)
	return r.Client.Publish(sub, data)
}

func (r *RpcNats) EncodeMsg(msgType MessageType, msgId int32, req interface{}) ([]byte, error) {
	rpcMsg := &RpcMsg{
		MsgType: msgType,
		MsgId:   msgId,
		MsgData: req,
	}
	data, err := r.RpcCoder.Encode(rpcMsg)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RpcNats) DecodeMsg(data []byte, v interface{}) error {
	return r.RpcCoder.DecodeMsg(data, v)
}

func (r *RpcNats) GetCoder() *RpcEncoder {
	return r.RpcCoder
}

func (r *RpcNats) Response(v interface{}) []byte {
	return r.RpcCoder.Response(v)
}

func (r *RpcNats) GetServer() *treaty.Server {
	return r.Server
}
