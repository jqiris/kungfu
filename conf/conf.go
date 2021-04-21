package conf

import (
	"github.com/jqiris/zinx/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var (
	config = new(Config)
	logger = logrus.WithField("package", "conf")
)

func InitConf(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("read file: %v error:%v", filename, err)
		return err
	}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		logger.Error("decode yaml error: %v", err)
		return err
	}
	logger.Warnf("the conf is:%+v", config)
	return nil
}

func GetDiscoverConf() DiscoverConf {
	return config.Discover
}

func GetRpcxConf() RpcxConf {
	return config.Rpcx
}

func GetStoresConf() StoresConf {
	return config.Stores
}

func GetCoderConf() CoderConf {
	return config.Coder
}

func GetConnectorConf() utils.GlobalObj {
	return config.Connector
}
