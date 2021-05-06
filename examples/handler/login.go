package handler

import (
	"github.com/jqiris/zinx/ziface"
	"github.com/jqiris/zinx/znet"
)

type LogingHandler struct {
	znet.BaseRouter
}

func (s *LogingHandler) Handle(request ziface.IRequest) {
	//先读取客户端的数据
	logger.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
}
