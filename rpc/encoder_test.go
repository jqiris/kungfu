package rpc

import (
	"fmt"
	"testing"
	"time"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/serialize"
	"github.com/jqiris/kungfu/v2/treaty"
)

func TestRpcEncoder(t *testing.T) {
	coder := NewRpcEncoder(serialize.NewJsonSerializer())
	eData, err := coder.Encode(&MsgRpc{
		MsgType: MsgTypePublish,
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
	//msg := &MsgRpc{MsgData: &treaty.LoginRequest{}}
	msg := &MsgRpc{}
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

func TestDuration(t *testing.T) {
	var a time.Duration
	fmt.Println(a > 0)
}
