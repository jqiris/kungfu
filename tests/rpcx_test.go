package tests

import (
	"fmt"
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
	"reflect"
	"testing"
)

//func TestRpc(t *testing.T) {
//	cfg := config.GetRpcXConf()
//	//gate
//	s1 := &treaty.Server{
//		ServerId:   "1001",
//		ServerType: "string_Balancer",
//		ServerName: "gate",
//		ServerIp:   "127.0.0.1",
//		ClientPort: 123,
//	}
//	w1 := rpcx.NewRpcBalancer(cfg)
//	if err := w1.Subscribe(s1, func(req []byte) []byte {
//		logger.Infof("gate received: %v", string(req))
//		return []byte(fmt.Sprintf("gate received: %v", string(req)))
//	}); err != nil {
//		logger.Errorf("gate err:%v", err)
//	}
//	if err := w1.SubscribeBalancer(func(req []byte) []byte {
//		logger.Infof("gate2 received: %v", string(req))
//		return []byte(fmt.Sprintf("gate2 received: %v", string(req)))
//	}); err != nil {
//		logger.Errorf("gate2 err:%v", err)
//	}
//	//connector
//	s2 := &treaty.Server{
//		ServerId:   "1002",
//		ServerType: "string_Connector",
//		ServerName: "connector",
//		ServerIp:   "127.0.0.1",
//		ClientPort: 456,
//	}
//	w2 := rpcx.NewRpcConnector(cfg)
//	if err := w2.Subscribe(s2, func(req []byte) []byte {
//		logger.Infof("connector received: %v", string(req))
//		return []byte(fmt.Sprintf("connector received: %v", string(req)))
//	}); err != nil {
//		logger.Errorf("connector err:%v", err)
//	}
//	if err := w2.SubscribeConnector(func(req []byte) []byte {
//		logger.Infof("connector2 received: %v", string(req))
//		return []byte(fmt.Sprintf("connector2 received: %v", string(req)))
//	}); err != nil {
//		logger.Errorf("connector2 err:%v", err)
//	}
//	//connector
//	s3 := &treaty.Server{
//		ServerId:   "1003",
//		ServerType: "string_Game",
//		ServerName: "game",
//		ServerIp:   "127.0.0.1",
//		ClientPort: 789,
//	}
//	w3 := rpcx.NewRpcServer(cfg)
//	if err := w3.Subscribe(s3, func(req []byte) []byte {
//		logger.Infof("server received: %v", string(req))
//		return []byte(fmt.Sprintf("server received: %v", string(req)))
//	}); err != nil {
//		logger.Errorf("server err:%v", err)
//	}
//	if err := w3.SubscribeServer(func(req []byte) []byte {
//		logger.Infof("server2 received: %v", string(req))
//		return []byte(fmt.Sprintf("server2 received: %v", string(req)))
//	}); err != nil {
//		logger.Errorf("server2 err:%v", err)
//	}
//
//	//s1 请求 s2
//	reply, err := w1.Request(s2, []byte("from gate"))
//	logger.Infof("s1=>s2, reply:%v, err:%v", string(reply), err)
//}

func TestServerByte(t *testing.T) {
	server_id := "backend_1001"
	server_byte := []byte(server_id)
	fmt.Println(len(server_byte))
}

func TestServerInterface(t *testing.T) {
	var server rpcx.RpcServer
	server = rpcx.NewRpcServer(config.GetRpcXConf(), &treaty.Server{})
	value := reflect.TypeOf(server)
	fmt.Println(value.Kind(), value)
}

type RpcHandler struct {
	InType  reflect.Type
	OutType reflect.Type
	Handler reflect.Value
}

func BackendLogin(server rpcx.RpcServer, req *treaty.LoginRequest) *treaty.LoginResponse {
	fmt.Printf("server is:%+v, req is %+v \n", server, req)
	return &treaty.LoginResponse{
		Msg: "login response",
	}
}

func BackendLogin2(server rpcx.RpcServer, req *treaty.LoginRequest) {
	fmt.Printf("server is:%+v, req is %+v \n", server, req)
}

func TestRpcHandler(t *testing.T) {
	var server rpcx.RpcServer
	server = rpcx.NewRpcServer(config.GetRpcXConf(), &treaty.Server{})
	funcValue := reflect.ValueOf(BackendLogin)
	ts := reflect.TypeOf(BackendLogin)
	fmt.Println(funcValue, ts, ts.NumIn(), ts.NumOut(), ts.In(1), ts.Out(0))
	in := reflect.New(ts.In(1).Elem()).Interface()
	srv := reflect.ValueOf(server)
	srvT := ts.In(0).String()
	//_, ok := srvT.(rpcx.RpcServer)
	fmt.Println("11111", srvT, "2222", srv.Kind(), "3333")
	args := []reflect.Value{srv, reflect.ValueOf(in)}
	res := funcValue.Call(args)
	fmt.Println("res is:", res[0].Interface().(*treaty.LoginResponse))
}

func TestRpcHandler2(t *testing.T) {
	var server rpcx.RpcServer
	server = rpcx.NewRpcServer(config.GetRpcXConf(), &treaty.Server{})
	funcValue := reflect.ValueOf(BackendLogin2)
	ts := reflect.TypeOf(BackendLogin2)
	fmt.Println(funcValue, ts.Kind(), ts.NumIn(), ts.NumOut(), ts.In(1))
	in := reflect.New(ts.In(1).Elem()).Interface()
	args := []reflect.Value{reflect.ValueOf(server), reflect.ValueOf(in)}
	res := funcValue.Call(args)
	fmt.Println("res is:", res)
}
