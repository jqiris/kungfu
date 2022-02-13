package balancer

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/golang/protobuf/proto"
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpc"
	"github.com/jqiris/kungfu/serialize"
	"github.com/jqiris/kungfu/session"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/kungfu/utils"
	"net/http"
	"net/url"
)

type BaseBalancer struct {
	ServerId              string
	Server                *treaty.Server
	Rpc                   rpc.ServerRpc
	ClientServer          *http.Server
	ClientCoder           serialize.Serializer
	EventJsonSelf         rpc.CallbackFunc //处理自己的json事件
	EventHandlerSelf      rpc.CallbackFunc //处理自己的事件
	EventHandlerBroadcast rpc.CallbackFunc //处理广播事件
}

func (b *BaseBalancer) HandleBalance(w http.ResponseWriter, r *http.Request) {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	serverType, uid := "", 0
	if err == nil {
		if len(queryForm["server_type"]) > 0 {
			serverType = queryForm["server_type"][0]
		}
		if len(queryForm["uid"]) > 0 {
			uid = utils.StringToInt(queryForm["uid"][0])
		}
	}
	if len(serverType) < 1 {
		res := &treaty.BalanceResult{
			Code: treaty.CodeType_CodeChooseBackendLogin,
		}
		b.WriteResponse(w, res)
		return
	}
	connector, err := b.Balance(r.RemoteAddr)
	if err != nil {
		res := &treaty.BalanceResult{
			Code: treaty.CodeType_CodeFailed,
		}
		b.WriteResponse(w, res)
		return
	}
	backend := discover.GetServerByType(serverType, r.RemoteAddr)
	var backendPre *treaty.Server
	sess := session.GetSession(int32(uid))
	if sess != nil {
		backendPre = sess.Backend
	}
	res := &treaty.BalanceResult{
		Code:       treaty.CodeType_CodeSuccess,
		Connector:  connector,
		Backend:    backend,
		BackendPre: backendPre,
	}
	b.WriteResponse(w, res)
}

func (b *BaseBalancer) WriteResponse(w http.ResponseWriter, msg proto.Message) {
	if v, e := b.ClientCoder.Marshal(msg); e == nil {
		if _, e2 := w.Write(v); e2 != nil {
			logger.Error(e2)
		}
	}
}
func (b *BaseBalancer) Init() {
	if b.Server == nil {
		panic("服务配置信息不能为空")
		return
	}
	//赋值id
	b.ServerId = b.Server.ServerId
	//init the rpc
	b.Rpc = rpc.NewRpcServer(config.GetRpcConf(), b.Server)
	//init the coder
	b.ClientCoder = serialize.NewProtoSerializer()
	//set the server
	b.ClientServer = &http.Server{Addr: fmt.Sprintf("%s:%d", b.Server.ServerIp, b.Server.ClientPort)}
	//handle the blance
	http.HandleFunc("/balance", b.HandleBalance)
	//run the server
	go func() {
		err := b.ClientServer.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
		}
	}()
	logger.Info("init the balancer:", b.ServerId)
}

func (b *BaseBalancer) AfterInit() {
	if b.Server == nil {
		panic("服务配置信息不能为空")
		return
	}
	if b.EventJsonSelf == nil {
		panic("EventJsonSelf不能为空")
		return
	}
	if b.EventHandlerSelf == nil {
		panic("EventHandlerSelf不能为空")
		return
	}
	if b.EventHandlerBroadcast == nil {
		panic("EventHandlerBroadcast不能为空")
		return
	}
	builder := rpc.NewSubscriberRpc(b.Server).SetCodeType(rpc.CodeTypeProto).SetCallback(func(req *rpc.MsgRpc) []byte {
		return b.EventHandlerSelf(req)
	})
	//Subscribe event
	if err := b.Rpc.Subscribe(builder.Build()); err != nil {
		logger.Error(err)
	}
	builder = builder.SetSuffix("json").SetCodeType(rpc.CodeTypeJson).SetCallback(func(req *rpc.MsgRpc) []byte {
		return b.EventJsonSelf(req)
	})
	//Subscribe event
	if err := b.Rpc.Subscribe(builder.Build()); err != nil {
		logger.Error(err)
	}
	builder = builder.SetSuffix(rpc.DefaultSuffix).SetCodeType(rpc.CodeTypeProto).SetCallback(func(req *rpc.MsgRpc) []byte {
		return b.EventHandlerBroadcast(req)
	})
	if err := b.Rpc.SubscribeBalancer(builder.Build()); err != nil {
		logger.Error(err)
	}
	//register the service
	if err := discover.Register(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *BaseBalancer) BeforeShutdown() {
	//unregister the service
	if err := discover.UnRegister(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *BaseBalancer) Shutdown() {
	if b.ClientServer != nil {
		if err := b.ClientServer.Close(); err != nil {
			logger.Error(err)
		}
	}
	logger.Info("stop the balancer:", b.ServerId)
}

func (b *BaseBalancer) Balance(remoteAddr string) (*treaty.Server, error) {
	if server := discover.GetServerByType("connector", remoteAddr); server != nil {
		return server, nil
	}

	return nil, errors.New("no suitable connector found")
}

func (b *BaseBalancer) GetServer() *treaty.Server {
	return b.Server
}
