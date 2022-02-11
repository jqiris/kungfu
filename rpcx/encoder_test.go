package rpcx

import (
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/serialize"
	"github.com/jqiris/kungfu/treaty"
	"testing"
)

func TestRpcEncoder(t *testing.T) {
	coder := NewRpcEncoder(serialize.NewJsonSerializer())
	eData, err := coder.Encode(&RpcMsg{
		MsgType: Publish,
		MsgId:   int32(treaty.MsgId_Msg_Login_Request),
		MsgData: &treaty.LoginRequest{
			Uid:      1001,
			Nickname: "jason",
			Token:    "dss",
			Backend:  nil,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	//msg := &RpcMsg{MsgData: &treaty.LoginRequest{}}
	msg := &RpcMsg{}
	err = coder.Decode(eData, msg)
	if err != nil {
		t.Fatal(err)
	}
	logger.Infof("the msg is: %#v", msg)
	req := &treaty.LoginRequest{}
	err = coder.DecodeMsg(msg.MsgData.([]byte), req)
	if err != nil {
		t.Fatal(err)
	}
	logger.Infof("the req is: %#v", req)
}
