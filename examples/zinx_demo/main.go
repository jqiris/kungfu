package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jqiris/kungfu/v2/logger"

	"github.com/jqiris/kungfu/v2/config"

	"github.com/jqiris/kungfu/v2/discover"
	"github.com/jqiris/kungfu/v2/launch"
	"github.com/jqiris/kungfu/v2/stores"
)

func main() {
	//init conf
	if err := config.InitConf("./examples/zinx_demo/config.json"); err != nil {
		logger.Fatal(err)
	}
	//init discover
	discover.InitDiscoverer(config.GetDiscoverConf())

	//init stores
	stores.InitStoreKeeper(config.GetStoresConf())

	//launch servers
	launch.Startup()
	sg := make(chan os.Signal, 1)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	select {
	case s := <-sg:
		logger.Info("server got shutdown signal", s)
	}
	launch.Shutdown()
}
