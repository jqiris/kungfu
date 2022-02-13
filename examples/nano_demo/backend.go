package main

import (
	"errors"
	"fmt"
	"github.com/jqiris/kungfu/base"
	"github.com/jqiris/kungfu/launch"
	"github.com/jqiris/kungfu/rpc"
	"time"

	"github.com/jqiris/kungfu/channel"
	"github.com/jqiris/kungfu/treaty"

	"github.com/jqiris/kungfu/logger"
)

type MyBackend struct {
	*base.ServerBase
	handler *rpc.Handler
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

func (g *MyBackend) HandleSelfEvent(req *rpc.MsgRpc) []byte {
	logger.Infof("MyBackend HandleSelfEvent received: %+v", req)
	resp, err := g.handler.DealMsg(rpc.CodeTypeProto, g.Rpc, req)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resp
}

func (g *MyBackend) HandleBroadcastEvent(req *rpc.MsgRpc) []byte {
	logger.Infof("MyBackend HandleBroadcastEvent received: %+v \n", req)
	return nil
}
func MyBackendCreator(s *treaty.Server) (rpc.ServerEntity, error) {
	if len(s.ServerId) < 1 {
		return nil, errors.New("服务器id不能为空")
	}
	server := &MyBackend{
		ServerBase: base.NewServerBase(s),
		connMap:    make(map[int32]*treaty.Server),
	}
	handler := rpc.NewHandler()
	handler.Register(int32(treaty.RpcMsgId_RpcMsgBackendLogin), server.BackendLogin)
	handler.Register(int32(treaty.RpcMsgId_RpcMsgBackendLogout), server.BackendOut)
	handler.Register(int32(treaty.RpcMsgId_RpcMsgChatTest), server.ChannelTest)
	server.handler = handler
	return server, nil
}

func init() {
	launch.RegisterCreator(rpc.Server, MyBackendCreator)
}
