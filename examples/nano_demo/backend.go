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
	conns map[int32]*treaty.Server
}

func (b *MyBackend) EventHandleSelf(server rpcx.RpcServer, req *rpcx.RpcMsg) []byte {
	fmt.Printf("MyBackend EventHandleSelf received: %+v \n", req)
	msgId, msgData := treaty.RpcMsgId(req.MsgId), req.MsgData.([]byte)
	switch msgId {
	case treaty.RpcMsgId_RpcMsgBackendLogin:
		//服务端登录
		resp := &treaty.LoginResponse{}
		msg := &treaty.LoginRequest{}
		if err := server.DecodeMsg(msgData, msg); err != nil {
			logger.Error(err)
			resp.Code = treaty.CodeType_CodeFailed
			resp.Msg = err.Error()
			return server.Response(resp)
		} else {
			//检查游戏通道是否建立
			ch := channel.GetChannel(b.Server, msg.Uid)
			if ch != nil {
				ch.ReconnectNum++
				ch.ReconnectTime = time.Now().Unix()
				if err = channel.SaveChannel(ch); err != nil {
					logger.Error(err)
				}
				resp.Code = treaty.CodeType_CodeSuccess
				resp.Msg = "登录成功"
				resp.Backend = b.Server
				b.conns[msg.Uid] = msg.Connector
				return server.Response(resp)
			}
			//游戏通道建立
			ch = &treaty.GameChannel{
				Uid:        msg.Uid,
				Connector:  msg.Connector,
				Backend:    b.Server,
				CreateTime: time.Now().Unix(),
			}
			if err = channel.SaveChannel(ch); err != nil {
				logger.Error(err)
			}
			resp.Code = treaty.CodeType_CodeSuccess
			resp.Msg = "登录成功"
			resp.Backend = b.Server
			b.conns[msg.Uid] = msg.Connector
			return server.Response(resp)
		}
	case treaty.RpcMsgId_RpcMsgBackendLogout:
		//服务端登出
		resp := &treaty.LogoutResponse{}
		msg := &treaty.LogoutRequest{}
		if err := server.DecodeMsg(msgData, msg); err != nil {
			logger.Error(err)
			resp.Code = treaty.CodeType_CodeFailed
			resp.Msg = err.Error()
			return server.Response(resp)
		} else {
			//游戏机制检查
			//销毁通道
			if err := channel.DestroyChannel(msg.Backend, msg.Uid); err != nil {
				logger.Error(err)
			}
			resp.Code = treaty.CodeType_CodeSuccess
			resp.Msg = "登出成功"
			return server.Response(resp)
		}
	case treaty.RpcMsgId_RpcMsgChatTest:
		rMsg := fmt.Sprintf("received chat msg:%+v", string(msgData))
		logger.Infof(rMsg)
		resp := &treaty.ChannelMsgResponse{
			Code:      0,
			Msg:       "success",
			MsgData:   rMsg,
			Connector: nil,
		}
		return server.Response(resp)
	}
	logger.Errorf("undfined message:%+v", req)
	return nil
}

func (b *MyBackend) EventHandleBroadcast(server rpcx.RpcServer, req *rpcx.RpcMsg) []byte {
	fmt.Printf("MyBackend EventHandleBroadcast received: %+v \n", req)
	return nil
}

func init() {
	srv := &MyBackend{conns: make(map[int32]*treaty.Server)}
	srv.SetServerId("backend_3001")
	srv.RegEventHandlerSelf(srv.EventHandleSelf)
	srv.RegEventHandlerBroadcast(srv.EventHandleBroadcast)
	launch.RegisterServer(srv)
}
