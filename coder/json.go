package coder

import "encoding/json"

type JsonCoder struct {
}

func NewJsonCoder() *JsonCoder {
	return &JsonCoder{}
}

func (j *JsonCoder) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JsonCoder) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
