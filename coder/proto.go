package coder

import (
	"errors"
	"github.com/golang/protobuf/proto"
)

type ProtoCoder struct {
}

func NewProtoCoder() *ProtoCoder {
	return &ProtoCoder{}
}

func (p *ProtoCoder) Marshal(v interface{}) ([]byte, error) {
	if v2, ok := v.(proto.Message); ok {
		return proto.Marshal(v2)
	}
	return nil, errors.New("marshal not proto message")
}

func (p *ProtoCoder) Unmarshal(data []byte, v interface{}) error {
	if v2, ok := v.(proto.Message); ok {
		return proto.Unmarshal(data, v2)
	}
	return errors.New("unmarshal not proto message")
}
