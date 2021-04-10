package rpcx

import "github.com/jqiris/kungfu/common"

type CallbackFunc func(resp []byte)

//rpc interface
type RpcBase interface {
	Subscribe(server common.Server, callback CallbackFunc) error //self Subscribe
	Publish(server common.Server, data []byte) error             //publish
	Request(server common.Server, data []byte) ([]byte, error)   //request
}

type RpcGate interface {
	RpcBase
	SubscribeGate(callback CallbackFunc) error //gate subscribe
	PublishConnector(data []byte) error        //connect publish
	PublishServer(data []byte) error           //server publish
}

type RpcConnector interface {
	RpcBase
	SubscribeConnector(callback CallbackFunc) error //connect subscribe
	PublishGate(data []byte) error                  //gate publish
	PublishServer(data []byte) error                //server publish
}

type RpcServer interface {
	RpcBase
	SubscribeServer(callback CallbackFunc) error //server subscribe
	PublishConnector(data []byte) error          //connect publish
	PublishServer(data []byte) error             //server publish
}
