package balancer

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/apex/log"
	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/stores"
	"github.com/jqiris/kungfu/treaty"
)

type BaseBalancer struct {
	ServerId              int32
	Server                *treaty.Server
	Rpcx                  rpcx.RpcBalancer
	ClientServer          *http.Server
	ClientCoder           coder.Coder
	EventHandlerSelf      func(req []byte) []byte //处理自己的事件
	EventHandlerBroadcast func(req []byte) []byte //处理广播事件
}

func (b *BaseBalancer) HandleBalance(w http.ResponseWriter, r *http.Request) {
	server, err := b.Balance(r.RemoteAddr)
	if err != nil {
		res := &treaty.BalanceResult{
			Code: treaty.CodeType_CodeFailed,
		}
		if v, e := b.ClientCoder.Marshal(res); e == nil {
			if _, e2 := w.Write(v); e2 != nil {
				logger.Error(e2)
			}
		}
		return
	}
	res := &treaty.BalanceResult{
		Code:   treaty.CodeType_CodeSuccess,
		Server: server,
	}
	if v, e := b.ClientCoder.Marshal(res); e == nil {
		if _, e2 := w.Write(v); e2 != nil {
			logger.Error(e2)
		}
		return
	}
}
func (b *BaseBalancer) Init() {
	//find the  server config
	if b.Server = helper.FindServerConfig(conf.GetBalancerConf(), b.GetServerId()); b.Server == nil {
		logger.Fatal("BaseBalancer can find the server config")
	}
	//init the rpcx
	b.Rpcx = rpcx.NewRpcBalancer(conf.GetRpcxConf())
	//set the server
	b.ClientServer = &http.Server{Addr: fmt.Sprintf(":%d", b.Server.ClientPort)}
	//handle the blance
	http.HandleFunc("/blance", b.HandleBalance)
	//run the server
	go func() {
		err := b.ClientServer.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
		}
	}()
	logger.Infoln("init the balancer:", b.ServerId)
}

func (b *BaseBalancer) AfterInit() {
	//Subscribe event
	if err := b.Rpcx.Subscribe(b.Server, func(req []byte) []byte {
		logger.Infof("BaseBalancer Subscribe received: %+v", req)
		return b.EventHandlerSelf(req)
	}); err != nil {
		logger.Error(err)
	}
	if err := b.Rpcx.SubscribeBalancer(func(req []byte) []byte {
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
}

func (b *BaseBalancer) Balance(remoteAddr string) (*treaty.Server, error) {
	//set the key
	key := "/user_connector/" + remoteAddr
	//find from store
	if res, err := stores.Get(key); err == nil && res != nil {
		if server, ok := res.(*treaty.Server); ok {
			return server, nil
		}
	}
	//find connector
	list := discover.FindServer(treaty.ServerType_Connector)
	if listLen := len(list); listLen > 0 {
		server := list[rand.Intn(listLen)]
		//store the server
		if err := stores.Set(key, server, 0); err != nil {
			logger.Error(err)
		}
		return server, nil
	}
	return nil, errors.New("no suitable connector found")
}

func (b *BaseBalancer) GetServer() *treaty.Server {
	return b.Server
}

func (b *BaseBalancer) RegEventHandlerSelf(handler func(req []byte) []byte) { //注册自己事件处理器
	b.EventHandlerSelf = handler
}

func (b *BaseBalancer) RegEventHandlerBroadcast(handler func(req []byte) []byte) { //注册广播事件处理器
	b.EventHandlerBroadcast = handler
}
func (b *BaseBalancer) SetServerId(serverId int32) {
	b.ServerId = serverId
}

func (b *BaseBalancer) GetServerId() int32 {
	return b.ServerId
}
