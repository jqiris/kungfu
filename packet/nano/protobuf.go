package nano

import (
	"encoding/json"
	"io/ioutil"
)

type ProtoNano struct {
	Client map[string]string `json:"client"`
	Server map[string]string `json:"server"`
}

func LoadProtobuf(filename string) (*ProtoNano, error) {
	protos := new(ProtoNano)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("read file: %v error:%v", filename, err)
		return nil, err
	}
	err = json.Unmarshal(content, protos)
	if err != nil {
		logger.Error("decode json error: %v", err)
		return nil, err
	}
	logger.Warnf("the proto is:%+v", protos)
	return protos, nil
}
