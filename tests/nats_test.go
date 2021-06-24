package tests

import (
	"fmt"
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
	"github.com/nats-io/nats.go"
	"strings"
	"testing"
)

func TestNatsEncoder(t *testing.T) {
	rpcConf := config.GetRpcXConf()
	url := strings.Join(rpcConf.Endpoints, ",")
	nc, _ := nats.Connect(url)
	c, _ := nats.NewEncodedConn(nc, rpcx.NATS_ENCODER)
	defer c.Close()
	c.Subscribe("/nats_test", func(resp *treaty.LoginRequest) {
		fmt.Println("sub resp:", resp)
	})
	c.Publish("/nats_test", &treaty.LoginRequest{
		Uid:      111,
		Nickname: "jason",
		Token:    "dfs",
		Backend:  nil,
	})
	select {}
}
