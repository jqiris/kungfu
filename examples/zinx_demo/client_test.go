package main

import (
	"fmt"
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/packet/zinx"
	"github.com/jqiris/kungfu/serialize"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/treaty"
)

func TestClientLogin(t *testing.T) {
	//根据balancer获取connector服务器
	resp, err := http.Get("http://127.0.0.1:8188/balance?server_type=backend")
	if err != nil {
		logger.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal(err)
	}
	coder := serialize.NewProtoSerializer()
	res := &treaty.BalanceResult{}
	if err = coder.Unmarshal(bytes, res); err != nil {
		logger.Error(err)
		return
	}
	if res.Code > 0 {
		logger.Fatal("no suitable connector find:", res.Code)
	}
	//根据balancer获得connector连接地址，并发送登录消息
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", res.Connector.ServerIp, res.Connector.ClientPort))
	if err != nil {
		logger.Fatalf("client start err, exit:%v", err)
	}
	//发送登录信息
	reqData, err := coder.Marshal(&treaty.LoginRequest{
		Uid:      1001,
		Nickname: "jason",
		Token:    "ce0da27df7150196625e48c843deb1f9",
		Backend:  res.Backend,
	})
	if err != nil {
		logger.Fatal(err)
	}
	msg, _ := zinx.Encode(&zinx.Message{
		int32(treaty.MsgId_Msg_Login_Request),
		reqData,
	})
	_, err = conn.Write(msg)
	if err != nil {
		logger.Println("write error err ", err)
		return
	}
	recMsg, err := zinx.ReadMsg(conn)
	if err != nil {
		fmt.Println("server unpack err:", err)
		return
	}

	//解析data数据
	respData := &treaty.LoginResponse{}
	if err = coder.Unmarshal(recMsg.Data, respData); err != nil {
		logger.Printf("login received err:%v", err)
	}
	logger.Infof("login result is:%+v", respData)
	//登录成功后尝试发送一次聊天数据
	send := &treaty.ChannelMsgRequest{
		Uid:     1001,
		MsgData: "hello chat",
	}

	data1, err1 := coder.Marshal(send)
	if err1 != nil {
		logger.Fatal(err1)
	}

	msg, _ = zinx.Encode(&zinx.Message{
		Id:   int32(treaty.MsgId_Msg_Channel_Request),
		Data: data1,
	})
	_, err = conn.Write(msg)
	if err != nil {
		logger.Println("write error err ", err)
		return
	}
	recMsg, err = zinx.ReadMsg(conn)
	if err != nil {
		fmt.Println("server unpack err:", err)
		return
	}
	logger.Infof("received chat resp:%+v", recMsg.Data)
}

func TestTokenCreate(t *testing.T) {
	uid, nickname := 1001, "jason"
	tokenkey := config.GetConnectorConf().TokenKey
	token := helper.Md5(fmt.Sprintf("%d|%s|%s", uid, nickname, tokenkey))
	fmt.Println(token)
}
