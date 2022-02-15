package main

import (
	"fmt"
	"github.com/jqiris/kungfu/v2/base"
	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/discover"
	"github.com/jqiris/kungfu/v2/launch"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/packet/nano"
	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/session"
	"github.com/jqiris/kungfu/v2/tcpface"
	"github.com/jqiris/kungfu/v2/treaty"
	"github.com/jqiris/kungfu/v2/utils"
)

type UserConnector struct {
	*base.ServerConnector
}

func NewUserConnector() *UserConnector {
	return &UserConnector{}
}

func (u *UserConnector) Login(s *session.Session, req *treaty.LoginRequest) error {
	logger.Infof("login received: %+v", req)
	//回复信息
	resp := &treaty.LoginResponse{}
	uid, nickname := req.Uid, req.Nickname
	tokenKey := config.GetConnectorConf().TokenKey
	token := utils.Md5(fmt.Sprintf("%d|%s|%s", uid, nickname, tokenKey))
	if req.Token != token {
		resp.Code = treaty.CodeType_CodeFailed
		resp.Msg = "token不正确"
		return s.Response(resp)
	}
	//必须加入一个服务器
	if req.Backend == nil {
		resp.Code = treaty.CodeType_CodeChooseBackendLogin
		resp.Msg = "请选择后端服务器进行登录"
		return s.Response(resp)
	}
	//与后端服务器建立连接
	backend := discover.GetServerById(req.Backend.ServerId)
	if backend == nil {
		//查找同类服务器
		backend = discover.GetServerByType(req.Backend.ServerType, s.RemoteAddr().String())
	}
	if backend == nil {
		resp.Code = treaty.CodeType_CodeCannotFindBackend
		resp.Msg = "找不到服务器"
		return s.Response(resp)
	}
	//后端服务器进行登录操作
	respBack := &treaty.LoginResponse{}
	b := rpc.NewReqBuilder(backend).SetMsgId(int32(treaty.RpcMsgId_RpcMsgBackendLogin)).SetReq(req).SetResp(respBack).Build()
	if err := u.Rpc.Request(b); err != nil {
		resp.Code = treaty.CodeType_CodeFailed
		resp.Msg = err.Error()
		return s.Response(resp)
	} else {
		if respBack.Code == treaty.CodeType_CodeSuccess {
			//成功绑定session
			err = s.Bind(int64(req.Uid))
			if err != nil {
				resp.Code = treaty.CodeType_CodeFailed
				resp.Msg = err.Error()
				return s.Response(resp)
			}
			//设置后端服务器
			s.Set("backend", req.Backend)
		}
		var a uint64 = 1<<55 - 1
		logger.Infof("a is:%v", a)
		respBack.TestInt = a
		return s.Response(respBack)
	}
}

func (u *UserConnector) ChannelMsg(s *session.Session, req *treaty.ChannelMsgRequest) error {
	logger.Infof("ChannelMsg received: %+v", req)
	//回复信息
	resp := &treaty.ChannelMsgResponse{}
	if s.UID() < 1 {
		resp.Code = treaty.CodeType_CodeNotLogin
		resp.Msg = "请先登录"
		return s.Response(resp)
	}
	backend, ok := s.Value("backend").(*treaty.Server)
	if !ok || backend == nil {
		resp.Code = treaty.CodeType_CodeNotLogin
		resp.Msg = "请先登录2"
		return s.Response(resp)
	}
	bResp := &treaty.ChannelMsgResponse{}
	b := rpc.NewReqBuilder(backend).SetMsgId(int32(treaty.RpcMsgId_RpcMsgChatTest)).SetReq(req).SetResp(bResp).Build()
	if err := u.Rpc.Request(b); err != nil {
		resp.Code = treaty.CodeType_CodeFailed
		resp.Msg = err.Error()
		return s.Response(resp)
	} else {
		return s.Response(bResp)
	}
}

func (u *UserConnector) HandleSelfEvent(req *rpc.MsgRpc) []byte {
	fmt.Printf("MyConnector HandleSelfEvent received: %+v \n", req)

	msgId, msgData := treaty.RpcMsgId(req.MsgId), req.MsgData.([]byte)
	switch msgId {
	case treaty.RpcMsgId_RpcMsgMultiLoginOut:
		//多端登录退出，向客户端发消息
		msg := &treaty.MultiLoginOut{}
		if err := u.Rpc.DecodeMsg(rpc.CodeTypeProto, msgData, msg); err != nil {
			logger.Error(err)
		} else {
			logger.Info(msg)
		}
	}
	return nil
}

func (u *UserConnector) HandleBroadcastEvent(req *rpc.MsgRpc) []byte {
	fmt.Printf("MyConnector HandleBroadcastEvent received: %+v \n", req)
	return nil
}

func UserConnectorCreator(s *treaty.Server) (rpc.ServerEntity, error) {
	server := &UserConnector{ServerConnector: base.NewServerConnector(s)}
	server.RouteHandler = func(s tcpface.IServer) {
		rs := s.GetMsgHandler()
		router := rs.(*nano.MsgHandle)
		err := router.Register(server)
		if err != nil {
			logger.Fatal(err)
		}
	}
	server.SelfEventHandler = server.HandleSelfEvent
	server.BroadcastEventHandler = server.HandleBroadcastEvent
	return server, nil
}

func init() {
	launch.RegisterCreator(rpc.Connector, UserConnectorCreator)
}
