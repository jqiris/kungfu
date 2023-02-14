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

package tests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jqiris/kungfu/v2/config"
	"github.com/nats-io/nats.go"
)

func TestChannelSub(t *testing.T) {
	rpcConf := config.GetRpcConf()
	url := strings.Join(rpcConf.Endpoints, ",")
	nc, _ := nats.Connect(url)
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()

	type person struct {
		Name    string
		Address string
		Age     int
	}

	recvCh := make(chan *person)
	ec.BindRecvChan("hello", recvCh)

	sendCh := make(chan *person)
	ec.BindSendChan("hello", sendCh)

	me := &person{Name: "derek", Age: 22, Address: "140 New Montgomery Street"}

	// Send via Go channels
	sendCh <- me

	// Receive via Go channels

	who := <-recvCh
	fmt.Println(who)
}

func TestChannelSubProto(t *testing.T) {
	rpcConf := config.GetRpcConf()
	url := strings.Join(rpcConf.Endpoints, ",")
	nc, _ := nats.Connect(url)
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()

	type person struct {
		Name    string
		Address string
		Age     int
	}
	type msg struct {
		Subject string
		Reply   string
		Person  *person
	}

	recvCh := make(chan *msg)
	sub, _ := ec.BindRecvChan("hello", recvCh)
	fmt.Println(sub)

	sendCh := make(chan *msg)
	ec.BindSendChan("hello", sendCh)

	//me := &treaty.RpcMsg{
	//	MsgId: treaty.RpcMsgId_RpcMsgBackendLogin,
	//	MsgServer: &treaty.Server{
	//		ServerId:   "999",
	//		ServerType: "sdsfs",
	//		ServerName: "sfds",
	//		ServerIp:   "dsss",
	//		ClientPort: 9999,
	//	},
	//	MsgData: []byte("hello world"),
	//}
	me := &msg{
		Subject: "",
		Reply:   "",
		Person: &person{
			Name:    "jason",
			Address: "shanghai",
			Age:     18,
		},
	}

	// Send via Go channels
	sendCh <- me

	// Receive via Go channels

	who := <-recvCh
	fmt.Println(who, "---", who.Reply)
}

func TestChannelNats(t *testing.T) {
	rpcConf := config.GetRpcConf()
	url := strings.Join(rpcConf.Endpoints, ",")
	nc, _ := nats.Connect(url, nats.Name("good"))
	fmt.Println(nc.HeadersSupported())

	nc.Subscribe("hello", func(msg *nats.Msg) {
		fmt.Println(string(msg.Data), msg.Header.Get("msgId"), msg.Header.Get("msgSource"))
	})

	req := nats.NewMsg("hello")
	req.Data = []byte("hello world")
	req.Header = make(map[string][]string)
	req.Header.Set("msgId", "1")
	req.Header.Set("msgSource", "server_001")
	resp, err := nc.RequestMsg(req, 10*time.Second)
	fmt.Println(resp, err)
}
