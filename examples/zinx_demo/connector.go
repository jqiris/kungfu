package main

import (
	"fmt"
	"github.com/jqiris/kungfu/packet/zinx"

	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/connector"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/launch"
	"github.com/jqiris/kungfu/session"
	"github.com/jqiris/kungfu/tcpface"
	"github.com/jqiris/kungfu/treaty"
)

type MyConnector struct {
	connector.TcpConnector
	conns map[int32]tcpface.IConnection
}

func (b *MyConnector) EventHandleSelf(req []byte) []byte {
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
			if conn, ok := b.conns[msg.Uid]; ok {
				SendMsg(conn, treaty.MsgId_Msg_Multi_Login_Out, msg)
				delete(b.conns, msg.Uid)
			}
		}
	}
	return nil
}

func (b *MyConnector) EventHandleBroadcast(req []byte) []byte {
	fmt.Printf("MyConnector EventHandleBroadcast received: %+v \n", string(req))
	return nil
}

//Login 登录操作
func (b *MyConnector) Login(request *zinx.Request) {

	//先读取客户端的数据
	logger.Println("Login recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetMsgData()))

	//回复信息
	resp := &treaty.LoginResponse{}
	//回复对象
	conn := request.GetConnection()
	//解析登录数据
	req := &treaty.LoginRequest{}
	if err := GetRequest(request, req); err != nil {
		resp.Code = treaty.CodeType_CodeFailed
		resp.Msg = err.Error()
		SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
		return
	}
	logger.Printf("login request is:%+v", req)
	//判断登录信息的正确性
	uid, nickname := req.Uid, req.Nickname
	tokenkey := config.GetConnectorConf().TokenKey
	token := helper.Md5(fmt.Sprintf("%d|%s|%s", uid, nickname, tokenkey))
	if req.Token != token {
		resp.Code = treaty.CodeType_CodeFailed
		resp.Msg = "token不正确"
		SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
		return
	}
	//必须加入一个服务器
	if req.Backend == nil {
		resp.Code = treaty.CodeType_CodeChooseBackendLogin
		resp.Msg = "请选择后端服务器进行登录"
		SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
		return
	}

	//检查会话信息
	sess := session.GetSession(uid)
	if sess != nil {
		if sess.Connector != nil && sess.Connector.ServerId != request.GetServerID() {
			//之前在其他客户端登录，通知其他connetor登出
			if msg, err := RpcMsgEncode(treaty.RpcMsgId_RpcMsgMultiLoginOut, &treaty.MultiLoginOut{Uid: uid}); err != nil {
				resp.Code = treaty.CodeType_CodeFailed
				resp.Msg = err.Error()
				SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
				return
			} else {
				if err = b.RpcX.Publish(sess.Connector, msg); err != nil {
					logger.Error(err)
				}
				//保存最新的connetor
				sess.Connector = discover.GetServerById(request.GetServerID())
				if err = session.SaveSession(uid, sess); err != nil {
					logger.Error(err)
				}
			}
		}
		//如果连接了后端服务器，进行重连
		if sess.Backend != nil && sess.Backend != req.Backend {
			resp.Code = treaty.CodeType_CodeLoginReconnect
			resp.Msg = "请进行重连服务器"
			resp.Backend = sess.Backend
			SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
			return
		}
	}

	//与后端服务器建立连接
	backend := discover.GetServerById(req.Backend.ServerId)
	if backend == nil {
		//查找同类服务器
		backend = discover.GetServerByType(req.Backend.ServerType, conn.RemoteAddr().String())
	}
	if backend == nil {
		resp.Code = treaty.CodeType_CodeCannotFindBackend
		resp.Msg = "找不到服务器"
		SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
		return
	}
	//后端服务器进行登录操作
	if msg, err := RpcMsgEncode(treaty.RpcMsgId_RpcMsgBackendLogin, req); err != nil {
		resp.Code = treaty.CodeType_CodeFailed
		resp.Msg = err.Error()
		SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
		return
	} else {
		if bResp, err := b.RpcX.Request(backend, msg); err != nil {
			resp.Code = treaty.CodeType_CodeFailed
			resp.Msg = err.Error()
			SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
			return
		} else {
			//结果直接由服务端返回
			respb := &treaty.LoginResponse{}
			if err = encoder.Unmarshal(bResp, respb); err != nil {
				resp.Code = treaty.CodeType_CodeFailed
				resp.Msg = err.Error()
				SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
				return
			}
			if resp.Code == treaty.CodeType_CodeSuccess {
				//登录成功记录用户的连接
				b.conns[uid] = conn
				//更新session
				if sess == nil {
					sess = &treaty.Session{Uid: uid}
				}
				sess.Connector = b.Server
				sess.Backend = backend
				if err := session.SaveSession(uid, sess); err != nil {
					logger.Error(err)
				}

			}
			SendMsg(conn, treaty.MsgId_Msg_Login_Response, respb)
			return
		}
	}
}

//Logout 登出操作
func (b *MyConnector) Logout(request *zinx.Request) {
	//先读取客户端的数据
	logger.Println("Logout recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetMsgData()))

	//回复信息
	resp := &treaty.LogoutResponse{}
	//回复对象
	conn := request.GetConnection()
	//解析登录数据
	req := &treaty.LogoutRequest{}
	if err := GetRequest(request, req); err != nil {
		resp.Code = treaty.CodeType_CodeFailed
		resp.Msg = err.Error()
		SendMsg(conn, treaty.MsgId_Msg_Logout_Response, resp)
		return
	}
	if req.Backend == nil {
		resp.Code = treaty.CodeType_CodeCannotFindBackend
		resp.Msg = "请指定登出服务器"
		SendMsg(conn, treaty.MsgId_Msg_Logout_Response, resp)
		return
	}
	logger.Printf("Logout request is:%+v", req)
	if msg, err := RpcMsgEncode(treaty.RpcMsgId_RpcMsgBackendLogout, req); err != nil {
		resp.Code = treaty.CodeType_CodeFailed
		resp.Msg = err.Error()
		SendMsg(conn, treaty.MsgId_Msg_Logout_Response, resp)
		return
	} else {
		if bResp, err := b.RpcX.Request(req.Backend, msg); err != nil {
			resp.Code = treaty.CodeType_CodeFailed
			resp.Msg = err.Error()
			SendMsg(conn, treaty.MsgId_Msg_Logout_Response, resp)
			return
		} else {
			//结果直接由服务端返回
			respb := &treaty.LogoutResponse{}
			if err = encoder.Unmarshal(bResp, respb); err != nil {
				resp.Code = treaty.CodeType_CodeFailed
				resp.Msg = err.Error()
				SendMsg(conn, treaty.MsgId_Msg_Logout_Response, resp)
				return
			}
			if resp.Code == treaty.CodeType_CodeSuccess {
				//登出成功删除用户连接
				delete(b.conns, req.Uid)
				//删除Session
				if err := session.DestorySession(req.Uid); err != nil {
					logger.Error(err)
				}
			}
			SendMsg(conn, treaty.MsgId_Msg_Logout_Response, respb)
			return
		}
	}
}

//ChannelMsg 消息转发
func (b *MyConnector) ChannelMsg(request *zinx.Request) {
	//先读取客户端的数据
	logger.Println("ChannelMsg recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetMsgData()))

	//回复信息
	resp := &treaty.ChannelMsgResponse{}
	//回复对象
	conn := request.GetConnection()
	//解析登录数据
	req := &treaty.ChannelMsgRequest{}
	if err := GetRequest(request, req); err != nil {
		resp.Code = treaty.CodeType_CodeFailed
		resp.Msg = err.Error()
		SendMsg(conn, treaty.MsgId_Msg_Channel_Response, resp)
		return
	}
	//检查session
	sess := session.GetSession(req.Uid)
	if sess == nil {
		resp.Code = treaty.CodeType_CodeNotLogin
		resp.Msg = "请先登录"
		SendMsg(conn, treaty.MsgId_Msg_Channel_Response, resp)
		return
	}
	if sess.Connector.ServerId != b.Server.ServerId {
		resp.Code = treaty.CodeType_CodeNotRightConnector
		resp.Msg = "请与绑定connector通信"
		SendMsg(conn, treaty.MsgId_Msg_Channel_Response, resp)
		return
	}
	if sess.Backend == nil {
		resp.Code = treaty.CodeType_CodeNotLoginBackend
		resp.Msg = "请先登录"
		SendMsg(conn, treaty.MsgId_Msg_Channel_Response, resp)
		return
	}
	if msg, err := encoder.Marshal(req.RpcMsg); err != nil {
		resp.Code = treaty.CodeType_CodeFailed
		resp.Msg = err.Error()
		SendMsg(conn, treaty.MsgId_Msg_Channel_Response, resp)
		return
	} else {
		if bResp, err := b.RpcX.Request(sess.Backend, msg); err != nil {
			resp.Code = treaty.CodeType_CodeFailed
			resp.Msg = err.Error()
			SendMsg(conn, treaty.MsgId_Msg_Channel_Response, resp)
			return
		} else {
			////结果直接由服务端返回
			//respb := &treaty.ChannelMsgResponse{}
			//if err = encoder.Unmarshal(bResp, respb); err != nil {
			//	resp.Code = treaty.CodeType_CodeFailed
			//	resp.Msg = err.Error()
			//	SendMsg(conn, treaty.MsgId_Msg_Channel_Response, resp)
			//	return
			//}
			SendByteMsg(conn, treaty.MsgId_Msg_Channel_Response, bResp)
			return
		}
	}
}

func init() {
	srv := &MyConnector{conns: make(map[int32]tcpface.IConnection)}
	srv.RouteHandler = func(s tcpface.IServer) {
		rs := s.GetMsgHandler()
		router := rs.(*zinx.MsgHandle)
		router.AddRouter(int32(treaty.MsgId_Msg_Login_Request), srv.Login)
		router.AddRouter(int32(treaty.MsgId_Msg_Channel_Request), srv.ChannelMsg)
	}
	srv.SetServerId("connector_2001")
	srv.RegEventHandlerSelf(srv.EventHandleSelf)
	srv.RegEventHandlerBroadcast(srv.EventHandleBroadcast)
	launch.RegisterServer(srv)
}
