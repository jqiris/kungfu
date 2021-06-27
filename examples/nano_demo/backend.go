package main

import (
	"fmt"
	"github.com/jqiris/kungfu/rpcx"
	"time"

	"github.com/jqiris/kungfu/channel"
	"github.com/jqiris/kungfu/treaty"

	"github.com/jqiris/kungfu/backend"
	"github.com/jqiris/kungfu/launch"
)

type MyBackend struct {
	backend.BaseBackEnd
	handler *rpcx.Handler
	conns   map[int32]*treaty.Server
}

func BackendLogin(server rpcx.RpcServer, req *treaty.LoginRequest) *treaty.LoginResponse {
	logger.Info("BackendLogin:", server, req)
	resp := &treaty.LoginResponse{}
	//检查游戏通道是否建立
	ch := channel.GetChannel(server.GetServer(), req.Uid)
	if ch != nil {
		ch.ReconnectNum++
		ch.ReconnectTime = time.Now().Unix()
		if err := channel.SaveChannel(ch); err != nil {
			logger.Error(err)
		}
		resp.Code = treaty.CodeType_CodeSuccess
		resp.Msg = "登录成功"
		resp.Backend = server.GetServer()
		return resp
	}
	//游戏通道建立
	ch = &treaty.GameChannel{
		Uid:        req.Uid,
		Connector:  req.Connector,
		Backend:    server.GetServer(),
		CreateTime: time.Now().Unix(),
	}
	if err := channel.SaveChannel(ch); err != nil {
		logger.Error(err)
	}
	resp.Code = treaty.CodeType_CodeSuccess
	resp.Msg = "登录成功"
	resp.Backend = server.GetServer()
	return resp
}

func BackendOut(server rpcx.RpcServer, msg *treaty.LogoutRequest) *treaty.LogoutResponse {
	//服务端登出
	resp := &treaty.LogoutResponse{}
	//销毁通道
	if err := channel.DestroyChannel(msg.Backend, msg.Uid); err != nil {
		logger.Error(err)
	}
	resp.Code = treaty.CodeType_CodeSuccess
	resp.Msg = "登出成功"
	return resp
}

func ChannelTest(server rpcx.RpcServer, req *treaty.ChannelMsgRequest) *treaty.ChannelMsgResponse {
	rMsg := fmt.Sprintf("received chat msg:%+v", string(req.MsgData))
	resp := &treaty.ChannelMsgResponse{
		Code:      0,
		Msg:       "success",
		MsgData:   rMsg,
		Connector: nil,
	}
	return resp
}

func (b *MyBackend) EventHandleSelf(server rpcx.RpcServer, req *rpcx.RpcMsg) []byte {
	fmt.Printf("MyBackend EventHandleSelf received: %+v \n", req)
	resp, err := b.handler.DealMsg(server, req)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resp
}

func (b *MyBackend) EventHandleBroadcast(server rpcx.RpcServer, req *rpcx.RpcMsg) []byte {
	fmt.Printf("MyBackend EventHandleBroadcast received: %+v \n", req)
	return nil
}

func init() {
	srv := &MyBackend{
		conns:   make(map[int32]*treaty.Server),
		handler: rpcx.NewHandler(),
	}
	srv.SetServerId("backend_3001")
	srv.RegEventHandlerSelf(srv.EventHandleSelf)
	srv.RegEventHandlerBroadcast(srv.EventHandleBroadcast)
	launch.RegisterServer(srv)
	srv.handler.Register(int32(treaty.RpcMsgId_RpcMsgBackendLogin), BackendLogin)
	srv.handler.Register(int32(treaty.RpcMsgId_RpcMsgBackendLogout), BackendOut)
	srv.handler.Register(int32(treaty.RpcMsgId_RpcMsgChatTest), ChannelTest)
}
