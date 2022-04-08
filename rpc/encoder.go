package rpc

import (
	"errors"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/serialize"
	"github.com/jqiris/kungfu/v2/utils"
)

type MessageType byte

// Message types
const (
	MsgTypeRequest  MessageType = 0x00
	MsgTypePublish              = 0x01
	MsgTypeResponse             = 0x02
)
const (
	msgHeadLength = 0x08
)

var (
	ErrInvalidMessage = errors.New("invalid message")
)

type EncoderRpc interface {
	Encode(rpcMsg *MsgRpc) ([]byte, error)
	Decode(data []byte, rpcMsg *MsgRpc) error
	DecodeMsg(data []byte, v any) error
	Response(v any) []byte
}

type MsgRpc struct {
	MsgType MessageType
	MsgId   int32
	MsgData any
}

type DefaultRpcEncoder struct {
	encoder serialize.Serializer
}

func NewRpcEncoder(encoder serialize.Serializer) *DefaultRpcEncoder {
	return &DefaultRpcEncoder{encoder: encoder}
}

// Encode Protocol
// --------<length>--------|--type--|----<MsgId>------|-<data>-
// ----------3byte---------|-1 byte-|-----4 byte------|--------
func (r *DefaultRpcEncoder) Encode(rpcMsg *MsgRpc) ([]byte, error) {
	var data []byte
	var err error
	switch rpcData := rpcMsg.MsgData.(type) {
	case string:
		data = []byte(rpcData)
	case []byte:
		data = rpcData
	default:
		data, err = r.encoder.Marshal(rpcMsg.MsgData)
		if err != nil {
			return nil, err
		}
	}
	//大端序
	length := msgHeadLength + len(data)
	buf := make([]byte, msgHeadLength)
	buf[0] = byte((length >> 16) & 0xFF)
	buf[1] = byte((length >> 8) & 0xFF)
	buf[2] = byte(length & 0xFF)
	buf[3] = byte(rpcMsg.MsgType)
	buf[4] = byte((rpcMsg.MsgId >> 24) & 0xFF)
	buf[5] = byte((rpcMsg.MsgId >> 16) & 0xFF)
	buf[6] = byte((rpcMsg.MsgId >> 8) & 0xFF)
	buf[7] = byte(rpcMsg.MsgId & 0xFF)
	buf = append(buf, data...)
	return buf, nil
}

func (r *DefaultRpcEncoder) Decode(data []byte, rpcMsg *MsgRpc) error {
	if len(data) < msgHeadLength {
		return ErrInvalidMessage
	}
	msgLength := utils.BigBytesToInt(data[:3])
	msgType := data[3]
	msgId := utils.BigBytesToInt(data[4:8])
	msgData := data[msgHeadLength:msgLength]
	rpcMsg.MsgType = MessageType(msgType)
	rpcMsg.MsgId = int32(msgId)
	if rpcMsg.MsgData == nil {
		rpcMsg.MsgData = msgData
	} else {
		err := r.encoder.Unmarshal(msgData, rpcMsg.MsgData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *DefaultRpcEncoder) DecodeMsg(data []byte, v any) error {
	err := r.encoder.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

func (r *DefaultRpcEncoder) Response(v any) []byte {
	rpcMsg := &MsgRpc{
		MsgType: MsgTypeResponse,
		MsgId:   0,
		MsgData: v,
	}
	res, err := r.Encode(rpcMsg)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return res
}
