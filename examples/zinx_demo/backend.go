package main

import (
	"fmt"
	"time"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/rpc"

	"github.com/jqiris/kungfu/v2/channel"
	"github.com/jqiris/kungfu/v2/treaty"

	"github.com/jqiris/kungfu/v2/launch"
)

type MyBackend struct {
	*rpc.ServerBase
	conns map[int32]*treaty.Server
}

func (b *MyBackend) HandleSelfEvent(req *rpc.MsgRpc) []byte {
	logger.Infof("MyBackend HandleSelfEvent received: %+v \n", req)
	msgId, msgData := treaty.RpcMsgId(req.MsgId), req.MsgData.([]byte)
	switch msgId {
	case treaty.RpcMsgId_RpcMsgBackendLogin:
		//服务端登录
		resp := &treaty.LoginResponse{}
		msg := &treaty.LoginRequest{}
		if err := b.Rpc.DecodeMsg(rpc.CodeTypeProto, msgData, msg); err != nil {
			logger.Error(err)
			resp.Code = treaty.CodeType_CodeFailed
			resp.Msg = err.Error()
			return b.Rpc.Response(rpc.CodeTypeProto, resp)
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
				return b.Rpc.Response(rpc.CodeTypeProto, resp)
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
			return b.Rpc.Response(rpc.CodeTypeProto, resp)
		}
	case treaty.RpcMsgId_RpcMsgBackendLogout:
		//服务端登出
		resp := &treaty.LogoutResponse{}
		msg := &treaty.LogoutRequest{}
		if err := b.Rpc.DecodeMsg(rpc.CodeTypeProto, msgData, msg); err != nil {
			logger.Error(err)
			resp.Code = treaty.CodeType_CodeFailed
			resp.Msg = err.Error()
			return b.Rpc.Response(rpc.CodeTypeProto, resp)
		} else {
			//游戏机制检查
			//销毁通道
			if err := channel.DestroyChannel(msg.Backend, msg.Uid); err != nil {
				logger.Error(err)
			}
			resp.Code = treaty.CodeType_CodeSuccess
			resp.Msg = "登出成功"
			return b.Rpc.Response(rpc.CodeTypeProto, resp)
		}
	case treaty.RpcMsgId_RpcMsgChatTest:
		resp := &treaty.ChannelMsgResponse{}
		msg := &treaty.ChannelMsgRequest{}
		if err := b.Rpc.DecodeMsg(rpc.CodeTypeProto, msgData, msg); err != nil {
			logger.Error(err)
			resp.Code = treaty.CodeType_CodeFailed
			resp.Msg = err.Error()
			return b.Rpc.Response(rpc.CodeTypeProto, resp)
		} else {
			resp.Code = 0
			resp.Msg = "success"
			resp.MsgData = fmt.Sprintf("received msg:%v", msg.MsgData)
			return b.Rpc.Response(rpc.CodeTypeProto, resp)
		}
	}
	logger.Errorf("undfined message:%+v", req)
	return nil
}

func MyBackendCreator(s *treaty.Server) (rpc.ServerEntity, error) {
	server := &MyBackend{
		ServerBase: rpc.NewServerBase(s),
		conns:      make(map[int32]*treaty.Server),
	}
	server.SetSelfEventHandler(server.HandleSelfEvent)
	return server, nil
}

func init() {
	launch.RegisterCreator("backend", MyBackendCreator)
}
