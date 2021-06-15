package zinx

import (
	"fmt"
	"github.com/jqiris/kungfu/utils"
)

// Message represents a unmarshaled message or a message which to be marshaled
type Message struct {
	Id   uint32 //消息的
	Data []byte //消息的内容
}

// NewMessage returns a new message instance
func NewMessage() *Message {
	return &Message{}
}

// String, implementation of fmt.Stringer interface
func (m *Message) String() string {
	return fmt.Sprintf("ID: %d,BodyLength: %d",
		m.Id,
		len(m.Data))
}

// Encode marshals message to binary format.
func (m *Message) Encode() ([]byte, error) {
	return MsgEncode(m)
}

// MsgEncode marshals message to binary format. Different message types is corresponding to
// different message header, message types is identified by 2-4 bit of flag field.
func MsgEncode(m *Message) ([]byte, error) {

	buf := make([]byte, 0)
	buf = append(buf, utils.LittleUInt32ToBytes(m.Id)...)
	buf = append(buf, m.Data...)
	return buf, nil
}

// MsgDecode unmarshal the bytes slice to a message
// See ref: https://github.com/lonnng/nano/blob/master/docs/communication_protocol.md
func MsgDecode(data []byte) (*Message, error) {
	m := NewMessage()
	offset := 4
	m.Id = utils.LittleBytesToUint32(data[:offset])
	m.Data = data[offset:]
	return m, nil
}
