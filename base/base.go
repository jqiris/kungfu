package base

import (
	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/discover"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/treaty"
)

type ServerBase struct {
	Server                *treaty.Server
	Rpc                   rpc.ServerRpc
	SubBuilder            *rpc.RssBuilder
	SelfEventHandler      rpc.CallbackFunc
	BroadcastEventHandler rpc.CallbackFunc
	handler               rpc.MsgHandler
}

func NewServerBase(s *treaty.Server) *ServerBase {
	return &ServerBase{
		Server:  s,
		handler: rpc.NewHandler(),
	}
}

func (s *ServerBase) SetMsgHandler(handler rpc.MsgHandler) {
	s.handler = handler
}

func (s *ServerBase) DealMsg(codeType string, server rpc.ServerRpc, req *rpc.MsgRpc) ([]byte, error) {
	return s.handler.DealMsg(codeType, server, req)
}

func (s *ServerBase) Register(msgId int32, v any) {
	s.handler.Register(msgId, v)
}

func (s *ServerBase) Init() {
	if s.Server == nil {
		panic("服务配置信息不能为空")
	}
	if len(s.Server.ServerId) < 1 || len(s.Server.ServerType) < 1 {
		panic("服务器基本配置信息不能为空")
	}
	//初始化rpc服务
	s.Rpc = rpc.NewRpcServer(config.GetRpcConf(), s.Server)
	//订阅创建
	s.SubBuilder = rpc.NewRssBuilder(s.Server)
	logger.Infof("init the service,type:%v, id:%v", s.Server.ServerType, s.Server.ServerId)
}

func (s *ServerBase) AfterInit() {
	if s.Server == nil {
		panic("服务配置信息不能为空")
	}
	if len(s.Server.ServerId) < 1 || len(s.Server.ServerType) < 1 {
		panic("服务器基本配置信息不能为空")
	}
	if s.SelfEventHandler == nil {
		panic("个体事件函数为空")
	}
	if s.BroadcastEventHandler == nil {
		panic("广播事件函数为空")
	}
	b := s.SubBuilder.Build()
	//sub self event
	if err := s.Rpc.Subscribe(b.SetCodeType(rpc.CodeTypeProto).SetCallback(s.SelfEventHandler).Build()); err != nil {
		panic(err)
	}
	//sub self json event
	b = s.SubBuilder.Build()
	if err := s.Rpc.Subscribe(b.SetSuffix(rpc.JsonSuffix).SetCodeType(rpc.CodeTypeJson).SetCallback(s.SelfEventHandler).Build()); err != nil {
		panic(err)
	}
	//sub broadcast event
	b = s.SubBuilder.Build()
	if err := s.Rpc.SubscribeBroadcast(b.SetCodeType(rpc.CodeTypeProto).SetCallback(s.BroadcastEventHandler).Build()); err != nil {
		panic(err)
	}
	//sub broadcast json event
	b = s.SubBuilder.Build()
	if err := s.Rpc.SubscribeBroadcast(b.SetSuffix(rpc.JsonSuffix).SetCodeType(rpc.CodeTypeJson).SetCallback(s.BroadcastEventHandler).Build()); err != nil {
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
