package config

import (
	"encoding/json"
	"github.com/jqiris/kungfu/v2/logger"
	"io/ioutil"

	"github.com/jqiris/kungfu/v2/treaty"
)

var (
	config = new(Config)
)

func InitConf(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("read file: %v error:%v", filename, err)
		return err
	}
	err = json.Unmarshal(content, config)
	if err != nil {
		logger.Error("decode json error: %v", err)
		return err
	}
	logger.Warnf("the config is:%+v", config)
	return nil
}

func InitFrameConf(content interface{}) error {
	bys, err := json.Marshal(content)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bys, config)
	if err != nil {
		logger.Error("decode json error: %v", err)
		return err
	}
	logger.Warnf("the config is:%+v", config)
	return nil
}

func GetDiscoverConf() DiscoverConf {
	return config.Discover
}

func GetRpcConf() RpcConf {
	return config.Rpc
}

func GetStoresConf() StoresConf {
	return config.Stores
}

func GetConnectorConf() ConnectorConf {
	return config.Connector
}

func GetServersConf() map[string]*treaty.Server {
	return config.Servers
}
