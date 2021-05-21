package tcpserver

//Message 消息
type Message struct {
	DataLen int32  //消息的长度
	ID      int32  //消息的ID
	Data    []byte //消息的内容
}

//NewMsgPackage 创建一个Message消息包
func NewMsgPackage(ID int32, data []byte) *Message {
	return &Message{
		DataLen: int32(len(data)),
		ID:      ID,
		Data:    data,
	}
}

//GetDataLen 获取消息数据段长度
func (msg *Message) GetDataLen() int32 {
	return msg.DataLen
}

//GetMsgID 获取消息ID
func (msg *Message) GetMsgID() int32 {
	return msg.ID
}

//GetData 获取消息内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

//SetDataLen 设置消息数据段长度
func (msg *Message) SetDataLen(len int32) {
	msg.DataLen = len
}

//SetMsgID 设计消息ID
func (msg *Message) SetMsgID(msgID int32) {
	msg.ID = msgID
}

//SetData 设计消息内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
