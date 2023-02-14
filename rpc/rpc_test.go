/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package rpc

import (
	"fmt"
	"testing"

	"github.com/jqiris/kungfu/v2/config"
)

var (
	testCfg = config.RpcConf{
		UseType:     "rabbitmq",
		DialTimeout: 5,
		Endpoints:   []string{"amqp://guest:guest@localhost:5672/mahjong"},
		DebugMsg:    true,
		Prefix:      "",
	}
)

func TestRabbitMqRprcSend(t *testing.T) {
	rrpc := NewRpcServer(testCfg, nil)
	defer rrpc.Close()
	req := DefaultReqBuilder().SetCodeType(CodeTypeJson).SetQueue("mj_queue").SetReq("模拟测试").SetMsgId(111).Build()
	err := rrpc.Publish(req)
	if err != nil {
		fmt.Println(err)
	}
}

func TestRabbitMqRprcReceive(t *testing.T) {
	rrpc := NewRpcServer(testCfg, nil)
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
