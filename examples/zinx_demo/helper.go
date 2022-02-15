package main

import (
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/packet/zinx"
	"github.com/jqiris/kungfu/v2/serialize"
	"github.com/jqiris/kungfu/v2/tcpface"
	"github.com/jqiris/kungfu/v2/treaty"
)

var (
	encoder = serialize.NewProtoSerializer()
)

func GetRequest(request *zinx.Request, v interface{}) error {
	if err := encoder.Unmarshal(request.GetMsgData(), v); err != nil {
		return err
	}
	return nil
}

func SendMsg(iConn tcpface.IConnection, msgId treaty.MsgId, msg interface{}) {
	conn := iConn.(*zinx.Agent)
	res, err := encoder.Marshal(msg)
	if err != nil {
		logger.Error(err)
		return
	}
	err = conn.SendBuffMsg(int32(msgId), res)
	if err != nil {
		logger.Error(err)
		return
	}
}
