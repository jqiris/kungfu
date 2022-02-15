package nano

import (
	"encoding/json"
	"github.com/jqiris/kungfu/v2/logger"
	"io/ioutil"
)

type ProtoNano struct {
	Client map[string]string `json:"client"`
	Server map[string]string `json:"server"`
}

func LoadProtobuf(filename string) (*ProtoNano, error) {
	ps := new(ProtoNano)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("read file: %v error:%v", filename, err)
		return nil, err
	}
	err = json.Unmarshal(content, ps)
	if err != nil {
		logger.Error("decode json error: %v", err)
		return nil, err
	}
	logger.Warnf("the proto is:%+v", ps)
	return ps, nil
}
