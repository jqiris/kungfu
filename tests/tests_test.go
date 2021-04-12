package tests

import (
	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/stores"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.WithField("package", "tests")
)

func init() {
	if err := conf.InitConf("../config.yaml"); err != nil {
		logger.Fatal(err)
	}
	//init discover
	discover.InitDiscoverer(conf.GetDiscoverConf())
	//init stores
	stores.InitStoreKeeper(conf.GetStoresConf())
	//init coder
	coder.InitCoder(conf.GetCoderConf())
}
