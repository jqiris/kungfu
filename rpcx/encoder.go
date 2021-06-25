package rpcx

import (
	"errors"
	"github.com/jqiris/kungfu/serialize"
	"github.com/jqiris/kungfu/utils"
)

const (
	NATS_ENCODER = "nats"
)

type NatsEncoder struct {
	useType string
	encoder serialize.Serializer
}

func NewNatsEncoder(useType string) *NatsEncoder {
	var encoder serialize.Serializer
	switch useType {
	case "json":
		encoder = serialize.NewJsonSerializer()
	case "proto":
		encoder = serialize.NewProtoSerializer()
	default:
		logger.Fatal("not support ")
	}
	return &NatsEncoder{
		useType: useType,
		encoder: encoder,
	}
}

func (n *NatsEncoder) Encode(subject string, v interface{}) ([]byte, error) {
	logger.Infof("nats encoder subject: %s", subject)
	return n.encoder.Marshal(v)
}

func (n *NatsEncoder) Decode(subject string, data []byte, vPtr interface{}) error {
	logger.Infof("nats decoder subject: %s", subject)
	return n.encoder.Unmarshal(data, vPtr)
}

type MessageType byte

// Message types
const (
	Request  MessageType = 0x00
	Notify               = 0x01
	Response             = 0x02
)
const (
	msgHeadLength = 0x06
)

var (
	ErrInvalidMessage = errors.New("invalid message")
)

type RpcMsg struct {
	MsgType MessageType
	MsgId   int32
	MsgData interface{}
}

type RpcEncoder struct {
	encoder serialize.Serializer
}

func NewRpcEncoder() *RpcEncoder {
	return &RpcEncoder{encoder: serialize.NewProtoSerializer()}
}

// Encode Protocol
// --------<length>--------|--type--|----<MsgId>------|-<data>-
// ----------3byte---------|-1 byte-|-----2 byte------|--------
func (r *RpcEncoder) Encode(rpcMsg *RpcMsg) ([]byte, error) {
	data, err := r.encoder.Marshal(rpcMsg.MsgData)
	if err != nil {
		return nil, err
	}
	//大端序
	length := msgHeadLength + len(data)
	buf := make([]byte, msgHeadLength)
	buf[0] = byte((length >> 16) & 0xFF)
	buf[1] = byte((length >> 8) & 0xFF)
	buf[2] = byte(length & 0xFF)
	buf[3] = byte(rpcMsg.MsgType)
	buf[4] = byte((rpcMsg.MsgId >> 8) & 0xFF)
	buf[5] = byte(rpcMsg.MsgId & 0xFF)
	buf = append(buf, data...)
	return buf, nil
}

func (r *RpcEncoder) Decode(data []byte, rpcMsg *RpcMsg) error {
	if len(data) < msgHeadLength {
		return ErrInvalidMessage
	}
	msgLength := utils.BigBytesToInt(data[:3])
	msgType := data[3]
	msgId := utils.BigBytesToInt(data[4:6])
	msgData := data[msgHeadLength:msgLength]
	rpcMsg.MsgType = MessageType(msgType)
	rpcMsg.MsgId = int32(msgId)
	err := r.encoder.Unmarshal(msgData, rpcMsg.MsgData)
	if err != nil {
		return err
	}
	return nil
}
