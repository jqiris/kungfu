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
	selfEventHandler      rpc.CallbackFunc
	broadcastEventHandler rpc.CallbackFunc
	innerMsgHandler       rpc.MsgHandler
	plugins               []rpc.ServerPlugin
}

func NewServerBase(s *treaty.Server, options ...Option) *ServerBase {
	server := &ServerBase{
		Server:          s,
		innerMsgHandler: rpc.NewHandler(),
		plugins:         make([]rpc.ServerPlugin, 0),
	}
	for _, option := range options {
		option(server)
	}
	return server
}

func (s *ServerBase) Register(msgId int32, v any) {
	s.innerMsgHandler.Register(msgId, v)
}

func (s *ServerBase) AddPlugin(plugin rpc.ServerPlugin) {
	s.plugins = append(s.plugins, plugin)
}

func (s *ServerBase) SetSelfEventHandler(handler rpc.CallbackFunc) {
	s.selfEventHandler = handler
}

func (s *ServerBase) SetBroadcastEventHandler(handler rpc.CallbackFunc) {
	s.broadcastEventHandler = handler
}

func (s *ServerBase) SetInnerMsgHandler(handler rpc.MsgHandler) {
	s.innerMsgHandler = handler
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
	//plugins
	for _, plugin := range s.plugins {
		plugin.Init(s.Server)
	}
	logger.Infof("init the service,type:%v, id:%v", s.Server.ServerType, s.Server.ServerId)
}

func (s *ServerBase) AfterInit() {
	if s.Server == nil {
		panic("服务配置信息不能为空")
	}
	if len(s.Server.ServerId) < 1 || len(s.Server.ServerType) < 1 {
		panic("服务器基本配置信息不能为空")
	}
	if s.selfEventHandler == nil {
		s.selfEventHandler = s.HandleSelfEvent
	}
	if s.broadcastEventHandler == nil {
		s.broadcastEventHandler = s.HandleBroadcastEvent
	}
	b := s.SubBuilder.Build()
	//sub self event
	if err := s.Rpc.Subscribe(b.SetCodeType(rpc.CodeTypeProto).SetCallback(s.selfEventHandler).Build()); err != nil {
		panic(err)
	}
	//sub self json event
	b = s.SubBuilder.Build()
	if err := s.Rpc.Subscribe(b.SetSuffix(rpc.JsonSuffix).SetCodeType(rpc.CodeTypeJson).SetCallback(s.selfEventHandler).Build()); err != nil {
		panic(err)
	}
	//sub broadcast event
	b = s.SubBuilder.Build()
	if err := s.Rpc.SubscribeBroadcast(b.SetCodeType(rpc.CodeTypeProto).SetCallback(s.broadcastEventHandler).Build()); err != nil {
		panic(err)
	}
	//sub broadcast json event
	b = s.SubBuilder.Build()
	if err := s.Rpc.SubscribeBroadcast(b.SetSuffix(rpc.JsonSuffix).SetCodeType(rpc.CodeTypeJson).SetCallback(s.broadcastEventHandler).Build()); err != nil {
		panic(err)
	}

	//plugins
	for _, plugin := range s.plugins {
		plugin.AfterInit(s.Server)
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
	//plugins
	for _, plugin := range s.plugins {
		plugin.BeforeShutdown()
	}
}

func (s *ServerBase) Shutdown() {
	//plugins
	for _, plugin := range s.plugins {
		plugin.Shutdown()
	}
	logger.Infof("shutdown service,type:%v,id:%v", s.Server.ServerType, s.Server.ServerId)
}

//内部事件处理
func (s *ServerBase) HandleSelfEvent(req *rpc.MsgRpc) []byte {
	resp, err := s.innerMsgHandler.DealMsg(rpc.CodeTypeProto, s.Rpc, req)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resp
}

func (s *ServerBase) HandleBroadcastEvent(req *rpc.MsgRpc) []byte {
	resp, err := s.innerMsgHandler.DealMsg(rpc.CodeTypeProto, s.Rpc, req)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resp
}
