package coder

import (
	"github.com/jqiris/kungfu/conf"
	"github.com/sirupsen/logrus"
)

var (
	logger   = logrus.WithField("package", "coder")
	defCoder Coder
)

func InitCoder(cfg conf.CoderConf) {
	switch cfg.UseType {
	case "json":
		defCoder = NewJsonCoder()
	case "proto":
		defCoder = NewProtoCoder()
	default:
		logger.Fatal("InitCoder failed")
	}
}

type Coder interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

func Marshal(v interface{}) ([]byte, error) {
	return defCoder.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return defCoder.Unmarshal(data, v)
}
