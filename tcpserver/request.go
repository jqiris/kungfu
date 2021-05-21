package tcpserver

//Request 请求
type Request struct {
	conn IConnection //已经和客户端建立好的 链接
	msg  IMessage    //客户端请求的数据
	sid  string      //客户端ID
}

//GetConnection 获取请求连接信息
func (r *Request) GetConnection() IConnection {
	return r.conn
}

//GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

//GetMsgID 获取请求的消息的ID
func (r *Request) GetMsgID() int32 {
	return r.msg.GetMsgID()
}

func (r *Request) GetServerID() string {
	return r.sid
}
