package main

import (
	"fmt"
	"github.com/jqiris/kungfu/rpcx"
	"time"

	"github.com/jqiris/kungfu/channel"
	"github.com/jqiris/kungfu/treaty"

	"github.com/jqiris/kungfu/backend"
	"github.com/jqiris/kungfu/launch"
	"github.com/jqiris/kungfu/logger"
)

type MyBackend struct {
	backend.BaseBackEnd
	handler *rpcx.Handler
	conns   map[int32]*treaty.Server
}

func (g *MyBackend) BackendLogin(req *treaty.LoginRequest) *treaty.LoginResponse {
	logger.Info("BackendLogin:", req)
	resp := &treaty.LoginResponse{}
	//检查游戏通道是否建立
	ch := channel.GetChannel(g.GetServer(), req.Uid)
	if ch != nil {
		ch.ReconnectNum++
		ch.ReconnectTime = time.Now().Unix()
		if err := channel.SaveChannel(ch); err != nil {
			logger.Error(err)
		}
		resp.Code = treaty.CodeType_CodeSuccess
		resp.Msg = "登录成功"
		resp.Backend = g.GetServer()
		return resp
	}
	//游戏通道建立
	ch = &treaty.GameChannel{
		Uid:        req.Uid,
		Connector:  req.Connector,
		Backend:    g.GetServer(),
		CreateTime: time.Now().Unix(),
	}
	if err := channel.SaveChannel(ch); err != nil {
		logger.Error(err)
	}
	resp.Code = treaty.CodeType_CodeSuccess
	resp.Msg = "登录成功"
	resp.Backend = g.GetServer()
	return resp
}

func (g *MyBackend) BackendOut(server rpcx.RpcServer, msg *treaty.LogoutRequest) *treaty.LogoutResponse {
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

func (g *MyBackend) ChannelTest(req *treaty.ChannelMsgRequest) *treaty.ChannelMsgResponse {
	rMsg := fmt.Sprintf("received chat msg:%+v", string(req.MsgData))
	resp := &treaty.ChannelMsgResponse{
		Code:      0,
		Msg:       "success",
		MsgData:   rMsg,
		Connector: nil,
	}
	return resp
}

func (b *MyBackend) EventHandleSelf(req *rpcx.RpcMsg) []byte {
	logger.Infof("MyBackend EventHandleSelf received: %+v", req)
	resp, err := b.handler.DealMsg(b.RpcX, req)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resp
}

func (b *MyBackend) EventHandleBroadcast(req *rpcx.RpcMsg) []byte {
	fmt.Printf("MyBackend EventHandleBroadcast received: %+v \n", req)
	return nil
}

func init() {
	srv := &MyBackend{
		conns:   make(map[int32]*treaty.Server),
		handler: rpcx.NewHandler(0, 0),
	}
	srv.SetServerId("backend_3001")
	srv.RegEventHandlerSelf(srv.EventHandleSelf)
	srv.RegEventHandlerBroadcast(srv.EventHandleBroadcast)
	launch.RegisterServer(srv)
	srv.handler.Register(int32(treaty.RpcMsgId_RpcMsgBackendLogin), srv.BackendLogin)
	srv.handler.Register(int32(treaty.RpcMsgId_RpcMsgBackendLogout), srv.BackendOut)
	srv.handler.Register(int32(treaty.RpcMsgId_RpcMsgChatTest), srv.ChannelTest)
}
