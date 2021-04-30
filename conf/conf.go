package conf

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jqiris/kungfu/treaty"
	"github.com/sirupsen/logrus"
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
	err = json.Unmarshal(content, config)
	if err != nil {
		logger.Error("decode json error: %v", err)
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

func GetConnectorConf() ConnectorConf {
	return config.Connector
}

func GetServersConf() map[string]*treaty.Server {
	return config.Servers
}

func GetLauchConf() []string {
	return config.Launch
}

func IsInLauch(serverId string) bool {
	for _, v := range config.Launch {
		if v == serverId {
			return true
		}
	}
	return false
}
