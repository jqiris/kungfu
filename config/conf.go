package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jqiris/kungfu/v2/logger"

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

func InitFrameConf(content any) error {
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

func SetDiscoverConf(cfg DiscoverConf) {
	config.Discover = cfg
}

func GetRpcConf() RpcConf {
	return config.Rpc
}

func SetRpcConf(cfg RpcConf) {
	config.Rpc = cfg
}

func GetStoresConf() StoresConf {
	return config.Stores
}

func SetStoresConf(cfg StoresConf) {
	config.Stores = cfg
}

func GetConnectorConf() ConnectorConf {
	return config.Connector
}

func SetConnectorConf(cfg ConnectorConf) {
	config.Connector = cfg
}

func GetServersConf() map[string]*treaty.Server {
	return config.Servers
}

func SetServersConf(cfg map[string]*treaty.Server) {
	config.Servers = cfg
}
