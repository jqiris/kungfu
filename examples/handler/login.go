package handler

import (
	"fmt"

	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/session"
	"github.com/jqiris/kungfu/tcpserver"
	"github.com/jqiris/kungfu/treaty"
)

type LogingHandler struct {
	tcpserver.BaseRouter
}

func (s *LogingHandler) Handle(request tcpserver.IRequest) {

	//先读取客户端的数据
	logger.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	//回复信息
	resp := &treaty.LoginResponse{}
	//回复对象
	conn := request.GetConnection()
	//解析登录数据
	loginRequest := &treaty.LoginRequest{}
	if err := GetRequest(request, loginRequest); err != nil {
		resp.Code = 1
		resp.Msg = err.Error()
		SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
	}
	logger.Printf("login request is:%+v", loginRequest)
	//判断登录信息的正确性
	uid, nickname := loginRequest.Uid, loginRequest.Nickname
	tokenkey := conf.GetConnectorConf().TokenKey
	token := helper.Md5(fmt.Sprintf("%d|%s|%s", uid, nickname, tokenkey))
	if loginRequest.Token != token {
		resp.Code = 1
		resp.Msg = "token不正确"
		SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
	}
	//检查会话信息
	sess := session.GetSession(uid)
	if sess != nil {
		if sess.Connector != nil && sess.Connector.ServerId != request.GetServerID() {
			//之前在其他客户端登录，通知其他connetor登出
		}
	}
	resp.Code = 0
	resp.Msg = "登录成功"
	SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
}
