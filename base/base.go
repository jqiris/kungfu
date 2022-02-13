package base

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpc"
	"github.com/jqiris/kungfu/treaty"
)

type ServerBase struct {
	Server     *treaty.Server
	Rpc        rpc.ServerRpc
	SubBuilder *rpc.SubscriberRpc
}

func NewServerBase(s *treaty.Server) *ServerBase {
	return &ServerBase{Server: s}
}

func (s *ServerBase) Init() {
	if s.Server == nil {
		panic("服务配置信息不能为空")
	}
	//初始化rpc服务
	s.Rpc = rpc.NewRpcServer(config.GetRpcConf(), s.Server)
	//订阅创建
	s.SubBuilder = rpc.NewSubscriberRpc(s.Server)
	logger.Infof("init the service,type:%v, id:%v", s.Server.ServerType, s.Server.ServerId)
}

func (s *ServerBase) AfterInit() {
	if s.Server == nil {
		panic("服务配置信息不能为空")
	}
	b := s.SubBuilder.Build()
	//sub self event
	if err := s.Rpc.Subscribe(b.SetCodeType(rpc.CodeTypeProto).SetCallback(s.HandleSelfEvent).Build()); err != nil {
		panic(err)
	}
	//sub broadcast event
	b = s.SubBuilder.Build()
	if err := s.Rpc.SubscribeBroadcast(b.SetCodeType(rpc.CodeTypeProto).SetCallback(s.HandleBroadcastEvent).Build()); err != nil {
		panic(err)
	}
	//服务注册
	if err := discover.Register(s.Server); err != nil {
		panic(err)
	}
}

func (s *ServerBase) BeforeShutdown() {
	//服务卸载
	if err := discover.UnRegister(s.Server); err != nil {
		logger.Error(err)
	}
}

func (s *ServerBase) Shutdown() {
	logger.Infof("shutdown service,type:%v,id:%v", s.Server.ServerType, s.Server.ServerId)
}

func (s *ServerBase) HandleSelfEvent(req *rpc.MsgRpc) []byte {
	//TODO implement me
	panic("implement HandleSelfEvent")
}

func (s *ServerBase) HandleBroadcastEvent(req *rpc.MsgRpc) []byte {
	//TODO implement me
	panic("implement HandleBroadcastEvent")
}
