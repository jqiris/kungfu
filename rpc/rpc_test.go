package rpc

import (
	"fmt"
	"testing"
)

func TestRabbitMqRprcSend(t *testing.T) {
	rrpc := NewRabbitMqRpc(WithRabbitMqEndpoints([]string{"amqp://guest:guest@localhost:5672/mahjong"}))
	defer rrpc.Close()
	req := DefaultReqBuilder().SetCodeType(CodeTypeJson).SetQueue("mj_queue").SetReq("哈哈").SetMsgId(111).Build()
	err := rrpc.Publish(req)
	if err != nil {
		fmt.Println(err)
	}
}

func TestRabbitMqRprcReceive(t *testing.T) {
	rrpc := NewRabbitMqRpc(WithRabbitMqEndpoints([]string{"amqp://guest:guest@localhost:5672/mahjong"}))
	defer rrpc.Close()
	receiver := NewRssBuilder(nil).SetCodeType(CodeTypeJson).SetQueue("mj_queue").SetCallback(rabbitCallBack).Build()
	err := rrpc.Subscribe(receiver)
	if err != nil {
		fmt.Println(err)
	}
}

func rabbitCallBack(req *MsgRpc) []byte {
	fmt.Println("rabbitCallBack:", string(req.MsgData.([]byte)))
	return nil
}
