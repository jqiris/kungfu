package tests

import (
	"fmt"
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
	"github.com/nats-io/nats.go"
	"strings"
	"testing"
	"time"
)

func TestNatsEncoder(t *testing.T) {
	rpcConf := config.GetRpcXConf()
	url := strings.Join(rpcConf.Endpoints, ",")
	nc, _ := nats.Connect(url)
	c, _ := nats.NewEncodedConn(nc, rpcx.NATS_ENCODER)
	defer c.Close()

	c.Subscribe("/nats_test", func(subj, reply string, req *treaty.LoginRequest) {
		fmt.Printf("sub subj:%v,reply:%v, req:%+v\n", subj, reply, req)
		c.Publish(reply, &treaty.LoginResponse{
			Code:    0,
			Msg:     "success",
			Backend: nil,
		})
	})
	//c.Publish("/nats_test", &treaty.LoginRequest{
	//	Uid:      111,
	//	Nickname: "jason",
	//	Token:    "dfs",
	//	Backend:  nil,
	//})
	resp := &treaty.LoginResponse{}
	err := c.Request("/nats_test", &treaty.LoginRequest{
		Uid:      111,
		Nickname: "jason",
		Token:    "dfs",
		Backend:  nil,
	}, resp, 10*time.Second)

	if err == nil {
		logger.Printf("sub resp:%+v", resp)
	} else {
		logger.Errorf("sub resp error: %v", err)
	}

	select {}
}

func TestChannelSub(t *testing.T) {
	rpcConf := config.GetRpcXConf()
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
	rpcConf := config.GetRpcXConf()
	url := strings.Join(rpcConf.Endpoints, ",")
	nc, _ := nats.Connect(url)
	ec, _ := nats.NewEncodedConn(nc, rpcx.NATS_ENCODER)
	defer ec.Close()

	//type person struct {
	//	Name    string
	//	Address string
	//	Age     int
	//}

	recvCh := make(chan *treaty.RpcMsg)
	ec.BindRecvChan("hello", recvCh)

	sendCh := make(chan *treaty.RpcMsg)
	ec.BindSendChan("hello", sendCh)

	me := &treaty.RpcMsg{
		MsgId: treaty.RpcMsgId_RpcMsgBackendLogin,
		MsgServer: &treaty.Server{
			ServerId:   "999",
			ServerType: "sdsfs",
			ServerName: "sfds",
			ServerIp:   "dsss",
			ClientPort: 9999,
		},
		MsgData: []byte("hello world"),
	}

	// Send via Go channels
	sendCh <- me

	// Receive via Go channels

	who := <-recvCh
	fmt.Println(who)
}

func TestChannelNats(t *testing.T) {
	rpcConf := config.GetRpcXConf()
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

//func (r *RpcNats) Request(server *treaty.Server, msgId int32, req, resp interface{}) error {
//	var msg *nats.Msg
//	var err error
//	rpcMsg := &RpcMsg{
//		MsgType: Request,
//		MsgId:   msgId,
//		MsgData: req,
//	}
//	data, err := r.rpcEncoder.Encode(rpcMsg)
//	if err != nil {
//		return err
//	}
//	if msg, err = r.Client.Request("/rpcx/"+treaty.RegSeverItem(server), data, r.DialTimeout); err == nil {
//		respMsg := &RpcMsg{MsgData: resp}
//		err = r.rpcEncoder.Decode(msg.Data, respMsg)
//		if err != nil {
//			return err
//		}
//	} else {
//		return err
//	}
//	return nil
//}
//
//func (r *RpcNats) Notify(server *treaty.Server, msgId int32, req interface{}) error {
//	var err error
//	rpcMsg := &RpcMsg{
//		MsgType: Notify,
//		MsgId:   msgId,
//		MsgData: req,
//	}
//	data, err := r.rpcEncoder.Encode(rpcMsg)
//	if err != nil {
//		return err
//	}
//	if _, err = r.Client.Request("/rpcx/"+treaty.RegSeverItem(server), data, r.DialTimeout); err != nil {
//		return err
//	}
//	return nil
//}
