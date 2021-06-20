package nano

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type ProtoItem struct {
	Option string `json:"option"`
	Type   string `json:"type"`
	Tag    int    `json:"tag"`
}

type ProtoNano struct {
	Client map[string]map[string]interface{} `json:"client"`
	Server map[string]map[string]interface{} `json:"server"`
}

func ParseObject(obj map[string]interface{}) map[string]interface{} {
	proto := make(map[string]interface{})
	nestProtos := make(map[string]interface{})
	tags := make(map[int]string)
	for name, tag := range obj {
		params := strings.Split(name, " ")
		switch params[0] {
		case "message":
			if len(params) != 2 {
				continue
			}
			nestProtos[params[1]] = ParseObject(tag.(map[string]interface{}))
		case "required":
			fallthrough
		case "optional":
			fallthrough
		case "repeated":
			if len(params) != 3 {
				continue
			}
			tagInt := int(tag.(float64))
			proto[params[2]] = ProtoItem{
				Option: params[0],
				Type:   params[1],
				Tag:    tagInt,
			}
			tags[tagInt] = params[2]
		}
	}
	proto["__messages"] = nestProtos
	proto["__tags"] = tags
	return proto
}

func LoadProtos(filename string) (map[string]interface{}, error) {
	res := make(map[string]interface{})
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
	clients := make(map[string]interface{})
	for k, v := range protos.Client {
		clients[k] = ParseObject(v)
	}
	servers := make(map[string]interface{})
	for k, v := range protos.Server {
		servers[k] = ParseObject(v)
	}
	res["client"] = clients
	res["server"] = servers
	return res, nil
}
