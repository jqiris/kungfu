package handler

import (
	"fmt"
	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/zinx/ziface"
	"github.com/jqiris/zinx/znet"
)

type LogingHandler struct {
	znet.BaseRouter
}

func (s *LogingHandler) Handle(request ziface.IRequest) {
	encoder := coder.NewProtoCoder()
	//先读取客户端的数据
	logger.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	//回复信息
	resp := &treaty.LoginResponse{}
	//回复对象
	conn := request.GetConnection()
	//解析登录数据
	loginRequest := &treaty.LoginRequest{}
	if err := encoder.Unmarshal(request.GetData(), loginRequest); err != nil {
		resp.Code = 1
		resp.Msg = err.Error()
		if res, err := encoder.Marshal(resp); err == nil {
			if err = conn.SendBuffMsg(uint32(treaty.MsgId_Msg_Login_Response), res); err != nil {
				logger.Error(err)
			}
		}

	}
	logger.Printf("login request is:%+v", loginRequest)
	//判断登录信息的正确性
	uid, nickname := loginRequest.Uid, loginRequest.Nickname
	tokenkey := conf.GetConnectorConf().TokenKey
	token := helper.Md5(fmt.Sprintf("%d|%s|%s", uid, nickname, tokenkey))
	if loginRequest.Token != token {
		resp.Code = 1
		resp.Msg = "token不正确"
		if res, err := encoder.Marshal(resp); err == nil {
			if err = conn.SendBuffMsg(uint32(treaty.MsgId_Msg_Login_Response), res); err != nil {
				logger.Error(err)
			}
		}
	}

	resp.Code = 0
	resp.Msg = "登录成功"
	if res, err := encoder.Marshal(resp); err == nil {
		if err = conn.SendBuffMsg(uint32(treaty.MsgId_Msg_Login_Response), res); err != nil {
			logger.Error(err)
		}
	}
}
