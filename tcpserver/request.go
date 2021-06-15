package tcpserver

import "github.com/jqiris/kungfu/tcpface"

type Request struct {
	conn tcpface.IConnection //已经和客户端建立好的 链接
	msg  tcpface.IMessage    //客户端请求的数据
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() tcpface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID 获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}

// GetServerID 获取ServerId
func (r *Request) GetServerID() string {
	return ""
}
