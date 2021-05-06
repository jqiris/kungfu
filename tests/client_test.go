package tests

import (
	"encoding/json"
	"fmt"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/zinx/znet"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
)

func TestClientLogin(t *testing.T) {
	//根据balancer获取connector服务器
	resp, err := http.Get("http://127.0.0.1:8188/balance")
	if err != nil {
		logger.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal(err)
	}
	res := &treaty.BalanceResult{}
	if err := json.Unmarshal(bytes, res); err != nil {
		logger.Fatal(err)
	}
	if res.Code > 0 {
		logger.Fatal("no suitable connector find")
	}
	//根据balancer获得connector连接地址，并发送登录消息
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", res.Server.ClientPort))
	if err != nil {
		logger.Fatal("client start err, exit!")
	}
	//发封包message消息
	dp := znet.NewDataPack()
	msg, _ := dp.Pack(znet.NewMsgPackage(1, []byte("Zinx Client loggin Message")))
	_, err = conn.Write(msg)
	if err != nil {
		logger.Println("write error err ", err)
		return
	}
}
