package connector

import (
	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/zinx/utils"
	"github.com/jqiris/zinx/ziface"
	"github.com/jqiris/zinx/znet"
)

type BaseConnector struct {
	ServerId              string
	Server                *treaty.Server
	Rpcx                  rpcx.RpcConnector
	ClientServer          ziface.IServer
	ClientCoder           coder.Coder
	ConnectorConf         utils.GlobalObj
	EventHandlerSelf      func(req []byte) []byte //处理自己的事件
	EventHandlerBroadcast func(req []byte) []byte //处理广播事件
	Routers               map[uint32]ziface.IRouter
}

func (b *BaseConnector) Init() {
	//find the  server config
	if serverConf := helper.FindServerConfig(conf.GetServersConf(), b.GetServerId()); serverConf == nil {
		logger.Fatal("BaseConnector can't find the server config")
	} else {
		connectorConf := conf.GetConnectorConf()
		b.Server = serverConf
		b.ConnectorConf = utils.GlobalObj{
			ServerId:         serverConf.ServerId,
			ServerType:       serverConf.ServerType,
			ServerName:       serverConf.ServerName,
			ServerIp:         serverConf.ServerIp,
			ClientPort:       int(serverConf.ClientPort),
			Version:          connectorConf.Version,
			MaxPacketSize:    connectorConf.MaxPacketSize,
			MaxConn:          connectorConf.MaxConn,
			WorkerPoolSize:   connectorConf.WorkerPoolSize,
			MaxWorkerTaskLen: connectorConf.MaxWorkerTaskLen,
			MaxMsgChanLen:    connectorConf.MaxMsgChanLen,
			LogDir:           connectorConf.LogDir,
			LogFile:          connectorConf.LogFile,
			LogDebugClose:    connectorConf.LogDebugClose,
		}
	}
	//init the rpcx
	b.Rpcx = rpcx.NewRpcConnector(conf.GetRpcxConf())
	//run the front server
	b.ClientServer = znet.NewServer(b.ConnectorConf)
	b.ClientServer.AddRouters(b.Routers)
	go b.ClientServer.Serve()

	logger.Infoln("init the connector:", b.ServerId)
}

func (b *BaseConnector) AfterInit() {
	//Subscribe event
	if err := b.Rpcx.Subscribe(b.Server, func(req []byte) []byte {
		logger.Infof("BaseConnector Subscribe received: %+v", req)
		return b.EventHandlerSelf(req)
	}); err != nil {
		logger.Error(err)
	}
	if err := b.Rpcx.SubscribeConnector(func(req []byte) []byte {
		logger.Infof("BaseConnector SubscribeConnector received: %+v", req)
		return b.EventHandlerBroadcast(req)
	}); err != nil {
		logger.Error(err)
	}
	//register the service
	if err := discover.Register(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *BaseConnector) BeforeShutdown() {
	//unregister the service
	if err := discover.UnRegister(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *BaseConnector) Shutdown() {
	//stop the server
	if b.ClientServer != nil {
		b.ClientServer.Stop()
	}
	logger.Infoln("stop the connector:", b.ServerId)
}

func (b *BaseConnector) GetServer() *treaty.Server {
	return b.Server
}

func (b *BaseConnector) RegEventHandlerSelf(handler func(req []byte) []byte) { //注册自己事件处理器
	b.EventHandlerSelf = handler
}

func (b *BaseConnector) RegEventHandlerBroadcast(handler func(req []byte) []byte) { //注册广播事件处理器
	b.EventHandlerBroadcast = handler
}
func (b *BaseConnector) SetServerId(serverId string) {
	b.ServerId = serverId
}

func (b *BaseConnector) GetServerId() string {
	return b.ServerId
}

//RegRouters 注册路由函数
func (b *BaseConnector) RegRouters(routers map[uint32]ziface.IRouter) {
	b.Routers = routers
}
