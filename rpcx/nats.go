package rpcx

import (
	"github.com/jqiris/kungfu/treaty"
	"github.com/nats-io/nats.go"
	"strings"
	"time"
)

type RpcNats struct {
	Endpoints   []string
	Options     []nats.Option
	Client      *nats.Conn
	DialTimeout time.Duration
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
	return r
}

func (r *RpcNats) Subscribe(server treaty.Server, callback CallbackFunc) error {
	if _, err := r.Client.Subscribe("/rpcx/"+server.RegId(), func(msg *nats.Msg) {
		resp := callback(msg.Data)
		if resp != nil {
			if err := msg.Respond(resp); err != nil {
				logger.Error(err)
			}
		}
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) Publish(server treaty.Server, data []byte) error {
	if err := r.Client.Publish("/rpcx/"+server.RegId(), data); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) Request(server treaty.Server, data []byte) ([]byte, error) {
	var msg *nats.Msg
	var err error
	if msg, err = r.Client.Request("/rpcx/"+server.RegId(), data, r.DialTimeout); err == nil {
		return msg.Data, nil
	}
	return nil, err
}

func (r *RpcNats) SubscribeGate(callback CallbackFunc) error {
	if _, err := r.Client.Subscribe("/rpcx/gate", func(msg *nats.Msg) {
		resp := callback(msg.Data)
		if resp != nil {
			if err := msg.Respond(resp); err != nil {
				logger.Error(err)
			}
		}
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) SubscribeConnector(callback CallbackFunc) error {
	if _, err := r.Client.Subscribe("/rpcx/connnector", func(msg *nats.Msg) {
		resp := callback(msg.Data)
		if resp != nil {
			if err := msg.Respond(resp); err != nil {
				logger.Error(err)
			}
		}
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) SubscribeServer(callback CallbackFunc) error {
	if _, err := r.Client.Subscribe("/rpcx/server", func(msg *nats.Msg) {
		resp := callback(msg.Data)
		if resp != nil {
			if err := msg.Respond(resp); err != nil {
				logger.Error(err)
			}
		}
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) PublishGate(data []byte) error {
	if err := r.Client.Publish("/rpcx/gate", data); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) PublishConnector(data []byte) error {
	if err := r.Client.Publish("/rpcx/connector", data); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) PublishServer(data []byte) error {
	if err := r.Client.Publish("/rpcx/server", data); err != nil {
		return err
	}
	return nil
}
