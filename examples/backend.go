package examples

import (
	"fmt"
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

func (b *MyBackend) EventHandleSelf(req []byte) []byte {
	fmt.Printf("MyBackend EventHandleSelf received: %+v \n", string(req))
	rpcMsg, err := RpcMsgDecode(req)
	if err != nil {
		logger.Error(err)
		return nil
	}
	switch rpcMsg.MsgId {
	case treaty.RpcMsgId_RpcMsgBackendLogin:
		//服务端登录
		resp := &treaty.LoginResponse{}
		msg := &treaty.LoginRequest{}
		if err := encoder.Unmarshal(rpcMsg.MsgData, msg); err != nil {
			logger.Error(err)
			resp.Code = treaty.CodeType_CodeFailed
			resp.Msg = err.Error()
			return RpcResponse(resp)
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
				b.conns[msg.Uid] = rpcMsg.MsgServer
				return RpcResponse(resp)
			}
			//游戏通道建立
			ch = &treaty.GameChannel{
				Uid:        msg.Uid,
				Connector:  rpcMsg.MsgServer,
				Backend:    b.Server,
				CreateTime: time.Now().Unix(),
			}
			if err = channel.SaveChannel(ch); err != nil {
				logger.Error(err)
			}
			resp.Code = treaty.CodeType_CodeSuccess
			resp.Msg = "登录成功"
			resp.Backend = b.Server
			b.conns[msg.Uid] = rpcMsg.MsgServer
			return RpcResponse(resp)
		}
	case treaty.RpcMsgId_RpcMsgBackendLogout:
		//服务端登出
		resp := &treaty.LogoutResponse{}
		msg := &treaty.LogoutRequest{}
		if err := encoder.Unmarshal(rpcMsg.MsgData, msg); err != nil {
			logger.Error(err)
			resp.Code = treaty.CodeType_CodeFailed
			resp.Msg = err.Error()
			return RpcResponse(resp)
		} else {
			//游戏机制检查
			//销毁通道
			if err := channel.DestroyChannel(msg.Backend, msg.Uid); err != nil {
				logger.Error(err)
			}
			resp.Code = treaty.CodeType_CodeSuccess
			resp.Msg = "登出成功"
			return RpcResponse(resp)
		}
	case treaty.RpcMsgId_RpcMsgChatTest:
		logger.Infof("received chat msg:%+v", string(rpcMsg.MsgData))
		return []byte("received the chat msg:" + string(rpcMsg.MsgData))
	}
	logger.Errorf("undfined message:%+v", rpcMsg)
	return nil
}

func (b *MyBackend) EventHandleBroadcast(req []byte) []byte {
	fmt.Printf("MyBackend EventHandleBroadcast received: %+v \n", string(req))
	return nil
}

func init() {
	srv := &MyBackend{conns: make(map[int32]*treaty.Server)}
	srv.SetServerId("backend_3001")
	srv.RegEventHandlerSelf(srv.EventHandleSelf)
	srv.RegEventHandlerBroadcast(srv.EventHandleBroadcast)
	launch.RegisterServer(srv)
}
