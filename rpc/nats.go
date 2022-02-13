package rpc

import (
	"fmt"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/serialize"
	"path"
	"strings"
	"time"

	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/kungfu/utils"
	"github.com/nats-io/nats.go"
)

const (
	Balancer  = "balancer"
	Connector = "connector"
	Server    = "backend"
	Database  = "database"
)
const (
	DefaultQueue  = "dq"
	DefaultSuffix = ""
	JsonSuffix    = "json"
)
const (
	CodeTypeJson  = "json"
	CodeTypeProto = "proto"
)

func DefaultCallback(req *MsgRpc) []byte {
	logger.Info("DefaultCallback")
	return nil
}

type SubscriberRpc struct {
	queue    string
	server   *treaty.Server
	callback CallbackFunc
	codeType string
	suffix   string
	parallel bool
}

func NewSubscriberRpc(server *treaty.Server) *SubscriberRpc {
	return &SubscriberRpc{
		queue:    DefaultQueue,
		server:   server,
		callback: DefaultCallback,
		codeType: CodeTypeProto,
		suffix:   DefaultSuffix,
		parallel: true,
	}
}

func (r *SubscriberRpc) SetQueue(queue string) *SubscriberRpc {
	r.queue = queue
	return r
}

func (r *SubscriberRpc) SetServer(server *treaty.Server) *SubscriberRpc {
	r.server = server
	return r
}
func (r *SubscriberRpc) SetCallback(callback CallbackFunc) *SubscriberRpc {
	r.callback = callback
	return r
}
func (r *SubscriberRpc) SetCodeType(codeType string) *SubscriberRpc {
	r.codeType = codeType
	return r
}
func (r *SubscriberRpc) SetSuffix(suffix string) *SubscriberRpc {
	r.suffix = suffix
	return r
}
func (r *SubscriberRpc) SetParallel(parallel bool) *SubscriberRpc {
	r.parallel = parallel
	return r
}

func (r *SubscriberRpc) Build() SubscriberRpc {
	return SubscriberRpc{
		queue:    r.queue,
		server:   r.server,
		callback: r.callback,
		codeType: r.codeType,
		suffix:   r.suffix,
		parallel: r.parallel,
	}
}

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

func (r *NatsRpc) Find(serverType string, arg int) *treaty.Server {
	return r.Finder.GetUserServer(serverType, arg)
}

func (r *NatsRpc) RemoveFindCache(arg int) {
	r.Finder.RemoveUserCache(arg)
}

func (r *NatsRpc) prepare(s SubscriberRpc) (EncoderRpc, error) {
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return nil, fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	return coder, nil
}

func (r *NatsRpc) Subscribe(s SubscriberRpc) error {
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

func (r *NatsRpc) QueueSubscribe(s SubscriberRpc) error {
	coder, err := r.prepare(s)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, treaty.RegSeverItem(s.server), s.suffix)
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

func (r *NatsRpc) SubscribeBroadcast(s SubscriberRpc) error {
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

func (r *NatsRpc) Request(codeType, suffix string, server *treaty.Server, msgId int32, req, resp interface{}) error {
	coder := r.RpcCoder[codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", codeType)
	}
	var msg *nats.Msg
	var err error
	var data []byte
	data, err = r.EncodeMsg(coder, Request, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, treaty.RegSeverItem(server), suffix)
	if msg, err = r.Client.Request(sub, data, r.DialTimeout); err == nil {
		respMsg := &MsgRpc{MsgData: resp}
		err = coder.Decode(msg.Data, respMsg)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (r *NatsRpc) Publish(codeType, suffix string, server *treaty.Server, msgId int32, req interface{}) error {
	coder := r.RpcCoder[codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", codeType)
	}
	data, err := r.EncodeMsg(coder, Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, treaty.RegSeverItem(server), suffix)
	if err = r.Client.Publish(sub, data); err != nil {
		return err
	}
	return nil
}

func (r *NatsRpc) PublishBalancer(codeType, suffix string, msgId int32, req interface{}) error {
	coder := r.RpcCoder[codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", codeType)
	}
	data, err := r.EncodeMsg(coder, Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, Balancer, suffix)
	return r.Client.Publish(sub, data)
}

func (r *NatsRpc) PublishConnector(codeType, suffix string, msgId int32, req interface{}) error {
	coder := r.RpcCoder[codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", codeType)
	}
	data, err := r.EncodeMsg(coder, Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, Connector, suffix)
	return r.Client.Publish(sub, data)
}

func (r *NatsRpc) PublishServer(codeType, suffix string, msgId int32, req interface{}) error {
	coder := r.RpcCoder[codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", codeType)
	}
	data, err := r.EncodeMsg(coder, Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, Server, suffix)
	return r.Client.Publish(sub, data)
}

func (r *NatsRpc) PublishDatabase(codeType, suffix string, msgId int32, req interface{}) error {
	coder := r.RpcCoder[codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", codeType)
	}
	data, err := r.EncodeMsg(coder, Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, Database, suffix)
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
