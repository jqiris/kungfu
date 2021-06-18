package main

import (
	"fmt"
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/connector"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/launch"
	"github.com/jqiris/kungfu/packet/nano"
	"github.com/jqiris/kungfu/session"
	"github.com/jqiris/kungfu/tcpface"
	"github.com/jqiris/kungfu/treaty"
)

type UserConnector struct {
	connector.TcpConnector
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
	token := helper.Md5(fmt.Sprintf("%d|%s|%s", uid, nickname, tokenKey))
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
	if msg, err := RpcMsgEncode(treaty.RpcMsgId_RpcMsgBackendLogin, req); err != nil {
		resp.Code = treaty.CodeType_CodeFailed
		resp.Msg = err.Error()
		return s.Response(resp)
	} else {
		if bResp, err := u.RpcX.Request(backend, msg); err != nil {
			resp.Code = treaty.CodeType_CodeFailed
			resp.Msg = err.Error()
			return s.Response(resp)
		} else {
			//结果直接由服务端返回
			respBack := &treaty.LoginResponse{}
			if err = encoder.Unmarshal(bResp, respBack); err != nil {
				resp.Code = treaty.CodeType_CodeFailed
				resp.Msg = err.Error()
				return s.Response(resp)
			}
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
			return s.Response(respBack)
		}
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
	if msg, err := encoder.Marshal(req.RpcMsg); err != nil {
		resp.Code = treaty.CodeType_CodeFailed
		resp.Msg = err.Error()
		return s.Response(resp)
	} else {
		if bResp, err := u.RpcX.Request(backend, msg); err != nil {
			resp.Code = treaty.CodeType_CodeFailed
			resp.Msg = err.Error()
			return s.Response(resp)
		} else {
			return s.Response(bResp)
		}
	}
}

func (u *UserConnector) EventHandleSelf(req []byte) []byte {
	fmt.Printf("MyConnector EventHandleSelf received: %+v \n", string(req))

	rpcMsg, err := RpcMsgDecode(req)
	if err != nil {
		logger.Error(err)
		return nil
	}
	switch rpcMsg.MsgId {
	case treaty.RpcMsgId_RpcMsgMultiLoginOut:
		//多端登录退出，向客户端发消息
		msg := &treaty.MultiLoginOut{}
		if err := encoder.Unmarshal(rpcMsg.MsgData, msg); err != nil {
			logger.Error(err)
		} else {
			logger.Println(msg)
		}
	}
	return nil
}

func (u *UserConnector) EventHandleBroadcast(req []byte) []byte {
	fmt.Printf("MyConnector EventHandleBroadcast received: %+v \n", string(req))
	return nil
}

func init() {
	srv := NewUserConnector()
	srv.RouteHandler = func(s tcpface.IServer) {
		rs := s.GetMsgHandler()
		router := rs.(*nano.MsgHandle)
		err := router.Register(srv)
		if err != nil {
			logger.Fatal(err)
		}
	}
	srv.SetServerId("connector_2001")
	srv.RegEventHandlerSelf(srv.EventHandleSelf)
	srv.RegEventHandlerBroadcast(srv.EventHandleBroadcast)
	launch.RegisterServer(srv)
}
