package tests

import (
	"fmt"
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpc"
	"github.com/jqiris/kungfu/stores"
	"github.com/jqiris/kungfu/treaty"
	"testing"
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

func TestSub(t *testing.T) {
	s := rpc.NewRssBuilder(&treaty.Server{ServerId: "test_001"})
	fmt.Println(s)
	s1 := s.Build()
	s2 := s1.SetSuffix("wel").Build()
	fmt.Println(s, s1, s2)
}
