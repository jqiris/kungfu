package zinx

import (
	"github.com/jqiris/kungfu/tcpface"
)

type Request struct {
	agent *Agent
	msg   *Message
}

func (r *Request) GetConnID() int {
	return r.agent.connId
}

func (r *Request) GetMsgID() int32 {
	return r.msg.Id
}

func (r *Request) GetMsgData() []byte {
	return r.msg.Data
}

func (r *Request) GetConnection() tcpface.IConnection {
	return r.agent
}

func (r *Request) GetServerID() string {
	return r.agent.server.GetServerID()
}
