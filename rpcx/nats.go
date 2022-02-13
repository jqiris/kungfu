package rpcx

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
)
const (
	CodeTypeJson  = "json"
	CodeTypeProto = "proto"
)

func DefaultCallback(req *RpcMsg) []byte {
	logger.Info("DefaultCallback")
	return nil
}

type RpcSubscriber struct {
	queue    string
	server   *treaty.Server
	callback CallbackFunc
	codeType string
	suffix   string
	parallel bool
}

func NewRpcSubscriber(server *treaty.Server) *RpcSubscriber {
	return &RpcSubscriber{
		queue:    DefaultQueue,
		server:   server,
		callback: DefaultCallback,
		codeType: CodeTypeProto,
		suffix:   DefaultSuffix,
		parallel: false,
	}
}

func (r *RpcSubscriber) SetQueue(queue string) *RpcSubscriber {
	r.queue = queue
	return r
}

func (r *RpcSubscriber) SetServer(server *treaty.Server) *RpcSubscriber {
	r.server = server
	return r
}
func (r *RpcSubscriber) SetCallback(callback CallbackFunc) *RpcSubscriber {
	r.callback = callback
	return r
}
func (r *RpcSubscriber) SetCodeType(codeType string) *RpcSubscriber {
	r.codeType = codeType
	return r
}
func (r *RpcSubscriber) SetSuffix(suffix string) *RpcSubscriber {
	r.suffix = suffix
	return r
}
func (r *RpcSubscriber) SetParallel(parallel bool) *RpcSubscriber {
	r.parallel = parallel
	return r
}

type RpcNats struct {
	Endpoints   []string
	Options     []nats.Option
	Client      *nats.Conn
	DialTimeout time.Duration
	RpcCoder    map[string]RpcEncoder
	Server      *treaty.Server
	DebugMsg    bool
	Prefix      string
	Finder      *discover.Finder
}
type RpcNatsOption func(r *RpcNats)

func WithNatsDebugMsg(debug bool) RpcNatsOption {
	return func(r *RpcNats) {
		r.DebugMsg = debug
	}
}
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
func WithNatsPrefix(prefix string) RpcNatsOption {
	return func(r *RpcNats) {
		r.Prefix = prefix
	}
}

func NewRpcNats(opts ...RpcNatsOption) *RpcNats {
	r := &RpcNats{
		Prefix: "RpcX",
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
	r.RpcCoder = map[string]RpcEncoder{
		CodeTypeProto: NewRpcEncoder(serialize.NewProtoSerializer()),
		CodeTypeJson:  NewRpcEncoder(serialize.NewJsonSerializer()),
	}
	r.Finder = discover.NewFinder()
	return r
}

func (r *RpcNats) RegEncoder(typ string, encoder RpcEncoder) {
	if _, ok := r.RpcCoder[typ]; !ok {
		r.RpcCoder[typ] = encoder
	} else {
		logger.Fatalf("encoder type has exist:%v", typ)
	}
}

func (r *RpcNats) Find(serverType string, arg int) *treaty.Server {
	return r.Finder.GetUserServer(serverType, arg)
}

func (r *RpcNats) RemoveFindCache(arg int) {
	r.Finder.RemoveUserCache(arg)
}

func (r *RpcNats) prepare(s RpcSubscriber) (RpcEncoder, error) {
	coder := r.RpcCoder[s.codeType]
	if coder == nil {
		return nil, fmt.Errorf("rpc coder not exist:%v", s.codeType)
	}
	return coder, nil
}

func (r *RpcNats) Subscribe(s RpcSubscriber) error {
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

func (r *RpcNats) QueueSubscribe(s RpcSubscriber) error {
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

func (r *RpcNats) SubscribeBalancer(s RpcSubscriber) error {
	coder, err := r.prepare(s)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, Balancer, s.suffix)
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

func (r *RpcNats) SubscribeConnector(s RpcSubscriber) error {
	coder, err := r.prepare(s)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, Connector, s.suffix)
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

func (r *RpcNats) SubscribeServer(s RpcSubscriber) error {
	coder, err := r.prepare(s)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, Server, s.suffix)
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

func (r *RpcNats) SubscribeDatabase(s RpcSubscriber) error {
	coder, err := r.prepare(s)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, Database, s.suffix)
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

func (r *RpcNats) DealMsg(msg *nats.Msg, callback CallbackFunc, coder RpcEncoder) {
	req := &RpcMsg{}
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

func (r *RpcNats) Request(codeType, suffix string, server *treaty.Server, msgId int32, req, resp interface{}) error {
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
		respMsg := &RpcMsg{MsgData: resp}
		err = coder.Decode(msg.Data, respMsg)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (r *RpcNats) Publish(codeType, suffix string, server *treaty.Server, msgId int32, req interface{}) error {
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

func (r *RpcNats) PublishBalancer(codeType, suffix string, msgId int32, req interface{}) error {
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

func (r *RpcNats) PublishConnector(codeType, suffix string, msgId int32, req interface{}) error {
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

func (r *RpcNats) PublishServer(codeType, suffix string, msgId int32, req interface{}) error {
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

func (r *RpcNats) PublishDatabase(codeType, suffix string, msgId int32, req interface{}) error {
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

func (r *RpcNats) EncodeMsg(coder RpcEncoder, msgType MessageType, msgId int32, req interface{}) ([]byte, error) {
	rpcMsg := &RpcMsg{
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

func (r *RpcNats) DecodeMsg(codeType string, data []byte, v interface{}) error {
	coder := r.RpcCoder[codeType]
	if coder == nil {
		return fmt.Errorf("rpc coder not exist:%v", codeType)
	}
	return coder.DecodeMsg(data, v)
}

func (r *RpcNats) GetCoder(codeType string) RpcEncoder {
	return r.RpcCoder[codeType]
}

func (r *RpcNats) Response(codeType string, v interface{}) []byte {
	coder := r.RpcCoder[codeType]
	if coder == nil {
		logger.Errorf("rpc coder not exist:%v", codeType)
		return nil
	}
	return coder.Response(v)
}

func (r *RpcNats) GetServer() *treaty.Server {
	return r.Server
}
