package handler

import (
	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/zinx/ziface"
	"github.com/sirupsen/logrus"
)

var (
	logger  = logrus.WithField("package", "handler")
	encoder = coder.NewProtoCoder()
)

func GetRequest(request ziface.IRequest, v interface{}) error {
	if err := encoder.Unmarshal(request.GetData(), v); err != nil {
		return err
	}
	return nil
}

func SendMsg(conn ziface.IConnection, msgId treaty.MsgId, msg interface{}) {
	res, err := encoder.Marshal(msg)
	if err != nil {
		logger.Error(err)
		return
	}
	err = conn.SendBuffMsg(uint32(msgId), res)
	if err != nil {
		logger.Error(err)
		return
	}
}
