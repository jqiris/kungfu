package main

import (
	"fmt"
	"time"

	"github.com/jqiris/kungfu/v2/launch"
	"github.com/jqiris/kungfu/v2/rpc"

	"github.com/jqiris/kungfu/v2/channel"
	"github.com/jqiris/kungfu/v2/treaty"

	"github.com/jqiris/kungfu/v2/logger"
)

type MyBackend struct {
	*rpc.ServerBase
	connMap map[int32]*treaty.Server
}

func (g *MyBackend) BackendLogin(req *treaty.LoginRequest) *treaty.LoginResponse {
	logger.Info("BackendLogin:", req)
	resp := &treaty.LoginResponse{}
	//检查游戏通道是否建立
	ch := channel.GetChannel(g.Server, req.Uid)
	if ch != nil {
		ch.ReconnectNum++
		ch.ReconnectTime = time.Now().Unix()
		if err := channel.SaveChannel(ch); err != nil {
			logger.Error(err)
		}
		resp.Code = treaty.CodeType_CodeSuccess
		resp.Msg = "登录成功"
		resp.Backend = g.Server
		return resp
	}
	//游戏通道建立
	ch = &treaty.GameChannel{
		Uid:        req.Uid,
		Connector:  req.Connector,
		Backend:    g.Server,
		CreateTime: time.Now().Unix(),
	}
	if err := channel.SaveChannel(ch); err != nil {
		logger.Error(err)
	}
	resp.Code = treaty.CodeType_CodeSuccess
	resp.Msg = "登录成功"
	resp.Backend = g.Server
	return resp
}

func (g *MyBackend) BackendOut(msg *treaty.LogoutRequest) *treaty.LogoutResponse {
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

func MyBackendCreator(s *treaty.Server) (rpc.ServerEntity, error) {
	server := &MyBackend{
		ServerBase: rpc.NewServerBase(s),
		connMap:    make(map[int32]*treaty.Server),
	}
	server.Register(int32(treaty.RpcMsgId_RpcMsgBackendLogin), server.BackendLogin)
	server.Register(int32(treaty.RpcMsgId_RpcMsgBackendLogout), server.BackendOut)
	server.Register(int32(treaty.RpcMsgId_RpcMsgChatTest), server.ChannelTest)
	return server, nil
}

func init() {
	launch.RegisterCreator(rpc.Server, MyBackendCreator)
}
