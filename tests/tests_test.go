package tests

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/stores"
)

func init() {
	if err := config.InitConf("../examples/nano_demo/config.json"); err != nil {
		logger.Fatal(err)
	}
	//init discover
	discover.InitDiscoverer(config.GetDiscoverConf())
	//init stores
	stores.InitStoreKeeper(config.GetStoresConf())
}
