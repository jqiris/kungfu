package backend

import (
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
)

//mt BackEnd
type BaseBackEnd struct {
	Server *treaty.Server
	Rpcx   rpcx.RpcServer
}

func (b *BaseBackEnd) Init() {
	//run the server

}

func (b *BaseBackEnd) AfterInit() {
	//Subscribe event
	if err := b.Rpcx.Subscribe(b.Server, func(req []byte) []byte {
		logger.Infof("BaseBackEnd Subscribe received: %+v", req)
		return nil
	}); err != nil {
		logger.Error(err)
	}
	if err := b.Rpcx.SubscribeServer(func(req []byte) []byte {
		logger.Infof("BaseBalancer SubscribeBalancer received: %+v", req)
		return nil
	}); err != nil {
		logger.Error(err)
	}
	//register the service
	if err := discover.Register(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *BaseBackEnd) BeforeShutdown() {
	//unregister the service
	if err := discover.UnRegister(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *BaseBackEnd) Shutdown() {
	//shutdown server
}
