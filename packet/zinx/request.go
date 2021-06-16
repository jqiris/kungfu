package zinx

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
