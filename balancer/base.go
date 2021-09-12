package balancer

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/golang/protobuf/proto"
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpcx"
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
	RpcX                  rpcx.RpcServer
	ClientServer          *http.Server
	ClientCoder           serialize.Serializer
	EventHandlerSelf      rpcx.CallbackFunc //处理自己的事件
	EventHandlerBroadcast rpcx.CallbackFunc //处理广播事件
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
	//find the  server config
	if b.Server = utils.FindServerConfig(config.GetServersConf(), b.GetServerId()); b.Server == nil {
		logger.Fatal("BaseBalancer can find the server config")
	}
	//init the rpcx
	b.RpcX = rpcx.NewRpcServer(config.GetRpcXConf(), b.Server)
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
	//Subscribe event
	if err := b.RpcX.Subscribe(b.Server, func(req *rpcx.RpcMsg) []byte {
		logger.Infof("BaseBalancer Subscribe received: %+v", req)
		return b.EventHandlerSelf(req)
	}); err != nil {
		logger.Error(err)
	}
	if err := b.RpcX.SubscribeBalancer(func(req *rpcx.RpcMsg) []byte {
		logger.Infof("BaseBalancer SubscribeBalancer received: %+v", req)
		return b.EventHandlerBroadcast(req)
	}); err != nil {
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

func (b *BaseBalancer) RegEventHandlerSelf(handler rpcx.CallbackFunc) { //注册自己事件处理器
	b.EventHandlerSelf = handler
}

func (b *BaseBalancer) RegEventHandlerBroadcast(handler rpcx.CallbackFunc) { //注册广播事件处理器
	b.EventHandlerBroadcast = handler
}
func (b *BaseBalancer) SetServerId(serverId string) {
	b.ServerId = serverId
}

func (b *BaseBalancer) GetServerId() string {
	return b.ServerId
}
