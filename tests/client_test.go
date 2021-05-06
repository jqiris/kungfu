package tests

import (
	"encoding/json"
	"fmt"
	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/zinx/znet"
	"io"
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
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", res.Server.ServerIp, res.Server.ClientPort))
	if err != nil {
		logger.Fatal("client start err, exit!")
	}
	//发送登录信息
	encoder := coder.NewProtoCoder()
	dp := znet.NewDataPack()
	data, err := encoder.Marshal(&treaty.LoginRequest{
		Uid:      1001,
		Nickname: "jason",
		Token:    "ce0da27df7150196625e48c843deb1f9",
	})
	if err != nil {
		logger.Fatal(err)
	}
	msg, _ := dp.Pack(znet.NewMsgPackage(uint32(treaty.MsgId_Msg_Login_Request), data))
	_, err = conn.Write(msg)
	if err != nil {
		logger.Println("write error err ", err)
		return
	}
	//先读出流中的head部分
	headData := make([]byte, dp.GetHeadLen())
	_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
	if err != nil {
		fmt.Println("read head error")
		return
	}
	//将headData字节流 拆包到msg中
	msgHead, err := dp.Unpack(headData)
	if err != nil {
		fmt.Println("server unpack err:", err)
		return
	}

	if msgHead.GetDataLen() > 0 {
		//msg 是有data数据的，需要再次读取data数据
		recMsg := msgHead.(*znet.Message)
		recMsg.Data = make([]byte, recMsg.GetDataLen())

		//根据dataLen从io中读取字节流
		_, err := io.ReadFull(conn, recMsg.Data)
		if err != nil {
			fmt.Println("server unpack data err:", err)
			return
		}
		//解析data数据
		data := &treaty.LoginResponse{}
		if err = coder.Unmarshal(recMsg.Data, data); err == nil {
			logger.Printf("login received data:%+v", data)
		}

	}

}

func TestTokenCreate(t *testing.T) {
	uid, nickname := 1001, "jason"
	tokenkey := conf.GetConnectorConf().TokenKey
	token := helper.Md5(fmt.Sprintf("%d|%s|%s", uid, nickname, tokenkey))
	fmt.Println(token)
}
