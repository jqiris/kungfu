package rpc

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/jqiris/kungfu/v2/discover"
	"github.com/jqiris/kungfu/v2/serialize"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/treaty"
	"github.com/jqiris/kungfu/v2/utils"
	"github.com/nats-io/nats.go"
)

type NatsRpc struct {
	Endpoints   []string
	Options     []nats.Option
	Client      *nats.Conn
	DialTimeout time.Duration
	RpcCoder    map[string]EncoderRpc
	Server      *treaty.Server
	DebugMsg    bool
	Prefix      string
	Finder      *discover.Finder
}
type NatsRpcOption func(r *NatsRpc)

func WithNatsDebugMsg(debug bool) NatsRpcOption {
	return func(r *NatsRpc) {
		r.DebugMsg = debug
	}
}
func WithNatsEndpoints(endpoints []string) NatsRpcOption {
	return func(r *NatsRpc) {
		r.Endpoints = endpoints
	}
}
func WithNatsDialTimeout(timeout time.Duration) NatsRpcOption {
	return func(r *NatsRpc) {
		r.DialTimeout = timeout
	}
}
func WithNatsServer(server *treaty.Server) NatsRpcOption {
	return func(r *NatsRpc) {
		r.Server = server
	}
}
func WithNatsOptions(opts ...nats.Option) NatsRpcOption {
	return func(r *NatsRpc) {
		r.Options = opts
	}
}
func WithNatsPrefix(prefix string) NatsRpcOption {
	return func(r *NatsRpc) {
		r.Prefix = prefix
	}
}

func NewRpcNats(opts ...NatsRpcOption) *NatsRpc {
	r := &NatsRpc{
		Prefix: "Rpc",
	}
	for _, opt := range opts {
		opt(r)
	}
	url := strings.Join(r.Endpoints, ",")
	conn, err := nats.Connect(url, r.Options...)
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

func (r *NatsRpc) RegEncoder(typ string, encoder EncoderRpc) {
	if _, ok := r.RpcCoder[typ]; !ok {
		r.RpcCoder[typ] = encoder
	} else {
		logger.Fatalf("encoder type has exist:%v", typ)
	}
}

func (r *NatsRpc) Find(serverType string, arg any) *treaty.Server {
	return r.Finder.GetUserServer(serverType, arg)
}

func (r *NatsRpc) RemoveFindCache(arg any) {
	r.Finder.RemoveUserCache(arg)
}

func (r *NatsRpc) prepare(s RssBuilder) (EncoderRpc, error) {
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return nil, fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	return coder, nil
}

func (r *NatsRpc) Subscribe(s RssBuilder) error {
	coder, err := r.prepare(s)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, treaty.RegSeverItem(s.server), s.suffix)
	if _, err = r.Client.Subscribe(sub, func(msg *nats.Msg) {
		if s.parallel {
			go utils.SafeRun(func() {
				r.DealMsg(msg, s.callback, coder)
			})
		} else {
			utils.SafeRun(func() {
				r.DealMsg(msg, s.callback, coder)
			})
		}
	}); err != nil {
		return err
	}
	return nil
}

func (r *NatsRpc) QueueSubscribe(s RssBuilder) error {
	coder, err := r.prepare(s)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, treaty.RegSeverQueue(s.server.ServerType, s.queue), s.suffix)
	if _, err = r.Client.QueueSubscribe(sub, s.queue, func(msg *nats.Msg) {
		if s.parallel {
			go utils.SafeRun(func() {
				r.DealMsg(msg, s.callback, coder)
			})
		} else {
			utils.SafeRun(func() {
				r.DealMsg(msg, s.callback, coder)
			})
		}
	}); err != nil {
		return err
	}
	return nil
}

func (r *NatsRpc) SubscribeBroadcast(s RssBuilder) error {
	coder, err := r.prepare(s)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, s.server.ServerType, s.suffix)
	if _, err = r.Client.Subscribe(sub, func(msg *nats.Msg) {
		if s.parallel {
			go utils.SafeRun(func() {
				r.DealMsg(msg, s.callback, coder)
			})
		} else {
			utils.SafeRun(func() {
				r.DealMsg(msg, s.callback, coder)
			})
		}
	}); err != nil {
		return err
	}
	return nil
}

func (r *NatsRpc) DealMsg(msg *nats.Msg, callback CallbackFunc, coder EncoderRpc) {
	req := &MsgRpc{}
	err := coder.Decode(msg.Data, req)
	if err != nil {
		logger.Error(err)
		return
	}
	resp := callback(req)
	if resp != nil {
		if err = msg.Respond(resp); err != nil {
			logger.Error(err)
		}
	}
	if r.DebugMsg {
		logger.Infof("DealMsg,msgType: %v, msgId: %v", req.MsgType, req.MsgId)
	}
}

func (r *NatsRpc) Request(s ReqBuilder) error {
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	var msg *nats.Msg
	var err error
	var data []byte
	data, err = r.EncodeMsg(coder, Request, s.msgId, s.req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, treaty.RegSeverItem(s.server), s.suffix)
	if msg, err = r.Client.Request(sub, data, r.DialTimeout); err == nil {
		respMsg := &MsgRpc{MsgData: s.resp}
		err = coder.Decode(msg.Data, respMsg)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
func (r *NatsRpc) QueueRequest(s ReqBuilder) error {
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	var msg *nats.Msg
	var err error
	var data []byte
	data, err = r.EncodeMsg(coder, Request, s.msgId, s.req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, treaty.RegSeverQueue(s.serverType, s.queue), s.suffix)
	if msg, err = r.Client.Request(sub, data, r.DialTimeout); err == nil {
		respMsg := &MsgRpc{MsgData: s.resp}
		err = coder.Decode(msg.Data, respMsg)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (r *NatsRpc) Publish(s ReqBuilder) error {
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	data, err := r.EncodeMsg(coder, Publish, s.msgId, s.req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, treaty.RegSeverItem(s.server), s.suffix)
	if err = r.Client.Publish(sub, data); err != nil {
		return err
	}
	return nil
}

func (r *NatsRpc) QueuePublish(s ReqBuilder) error {
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	data, err := r.EncodeMsg(coder, Publish, s.msgId, s.req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, treaty.RegSeverQueue(s.serverType, s.queue), s.suffix)
	if err = r.Client.Publish(sub, data); err != nil {
		return err
	}
	return nil
}

func (r *NatsRpc) PublishBroadcast(s ReqBuilder) error {
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	data, err := r.EncodeMsg(coder, Publish, s.msgId, s.req)
	if err != nil {
		return err
	}
	if len(s.serverType) < 1 && s.server != nil {
		s.serverType = s.server.ServerType
	}
	sub := path.Join(r.Prefix, s.serverType, s.suffix)
	return r.Client.Publish(sub, data)
}

func (r *NatsRpc) EncodeMsg(coder EncoderRpc, msgType MessageType, msgId int32, req interface{}) ([]byte, error) {
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

func (r *NatsRpc) DecodeMsg(codeType string, data []byte, v interface{}) error {
	coder := r.RpcCoder[codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", codeType)
	}
	return coder.DecodeMsg(data, v)
}

func (r *NatsRpc) GetCoder(codeType string) EncoderRpc {
	return r.RpcCoder[codeType]
}

func (r *NatsRpc) Response(codeType string, v interface{}) []byte {
	coder := r.RpcCoder[codeType]
	if coder == nil {
		logger.Errorf("rpc coder not exist:%v", codeType)
		return nil
	}
	return coder.Response(v)
}

func (r *NatsRpc) GetServer() *treaty.Server {
	return r.Server
}
