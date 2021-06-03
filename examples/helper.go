package examples

import (
	"github.com/golang/protobuf/proto"
	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/tcpserver"
	"github.com/jqiris/kungfu/treaty"
	"github.com/sirupsen/logrus"
)

var (
	logger  = logrus.WithField("package", "handler")
	encoder = coder.NewProtoCoder()
)

func GetRequest(request tcpserver.IRequest, v interface{}) error {
	if err := encoder.Unmarshal(request.GetData(), v); err != nil {
		return err
	}
	return nil
}

func SendMsg(conn tcpserver.IConnection, msgId treaty.MsgId, msg interface{}) {
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
func SendByteMsg(conn tcpserver.IConnection, msgId treaty.MsgId, msg []byte) {
	err := conn.SendBuffMsg(int32(msgId), msg)
	if err != nil {
		logger.Error(err)
		return
	}
}

func RpcMsgEncode(msgId treaty.RpcMsgId, msgData proto.Message) ([]byte, error) {
	msg, err := encoder.Marshal(msgData)
	if err != nil {
		return nil, err
	}
	rpcMsg := &treaty.RpcMsg{
		MsgId:   msgId,
		MsgData: msg,
	}
	rpcLoad, err := encoder.Marshal(rpcMsg)
	if err != nil {
		return nil, err
	}
	return rpcLoad, nil
}

func RpcMsgDecode(data []byte) (*treaty.RpcMsg, error) {
	msg := &treaty.RpcMsg{}
	err := encoder.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func RpcResponse(msg proto.Message) []byte {
	if res, err := encoder.Marshal(msg); err != nil {
		logger.Error(err)
		return nil
	} else {
		return res
	}
}
