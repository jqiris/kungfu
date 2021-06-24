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
	if err := c.Request("/nats_test", &treaty.LoginRequest{
		Uid:      111,
		Nickname: "jason",
		Token:    "dfs",
		Backend:  nil,
	}, resp, 10*time.Second); err == nil {
		fmt.Printf("sub resp:%v", resp)
	} else {
		fmt.Errorf("sub resp error: %v", err)
	}

	select {}
}
