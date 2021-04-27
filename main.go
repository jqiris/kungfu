package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/discover"
	_ "github.com/jqiris/kungfu/examples"
	"github.com/jqiris/kungfu/launch"
	"github.com/jqiris/kungfu/stores"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.WithField("package", "main")
)

func main() {
	//init conf
	if err := conf.InitConf("config.json"); err != nil {
		logger.Fatal(err)
	}
	//init discover
	discover.InitDiscoverer(conf.GetDiscoverConf())

	//init stores
	stores.InitStoreKeeper(conf.GetStoresConf())

	//init coder
	coder.InitCoder(conf.GetCoderConf())

	//launch servers
	launch.Startup()
	sg := make(chan os.Signal, 1)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	select {
	case s := <-sg:
		log.Println("server got shutdown signal", s)
	}
	launch.Shutdown()
}
