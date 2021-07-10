package main

import (
	"github.com/jqiris/kungfu/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jqiris/kungfu/config"

	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/launch"
	"github.com/jqiris/kungfu/stores"
)

func main() {
	//run client
	//go RunClient()
	//init conf
	if err := config.InitConf("./examples/nano_demo/config.json"); err != nil {
		logger.Fatal(err)
	}
	//init discover
	discover.InitDiscoverer(config.GetDiscoverConf())

	//init stores
	stores.InitStoreKeeper(config.GetStoresConf())

	//init logger
	lg, cancel := logger.NewLogger(
		logger.WithOutType("out_all"),
		logger.WithLogDir("./logs"),
		logger.WithLogLevel("debug"),
		logger.WithLogName("nano_demo"),
		logger.WithLogRuntime(true),
		logger.WithLogDump(true),
		logger.WithStdColor(true),
		logger.WithZipDuration(5*time.Second),
		logger.WithTickTime(5*time.Second),
	)
	logger.SetLogger(lg, cancel)

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
