package rpc

import (
	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/discover"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/treaty"
)

type ServerBase struct {
	Server                *treaty.Server
	Rpc                   ServerRpc
	SubBuilder            *RssBuilder
	selfEventHandler      CallbackFunc
	broadcastEventHandler CallbackFunc
	innerMsgHandler       MsgHandler
	plugins               []ServerPlugin
}

func NewServerBase(s *treaty.Server, options ...Option) *ServerBase {
	server := &ServerBase{
		Server:          s,
		innerMsgHandler: NewHandler(),
		plugins:         make([]ServerPlugin, 0),
	}
	for _, option := range options {
		option(server)
	}
	return server
}

func (s *ServerBase) Register(msgId int32, v any) {
	s.innerMsgHandler.Register(msgId, v)
}

func (s *ServerBase) AddPlugin(plugin ServerPlugin) {
	s.plugins = append(s.plugins, plugin)
}

func (s *ServerBase) SetSelfEventHandler(handler CallbackFunc) {
	s.selfEventHandler = handler
}

func (s *ServerBase) SetBroadcastEventHandler(handler CallbackFunc) {
	s.broadcastEventHandler = handler
}

func (s *ServerBase) SetInnerMsgHandler(handler MsgHandler) {
	s.innerMsgHandler = handler
}

func (s *ServerBase) Init() {
	//init default rpc
	defRpcInit()
	//init current server
	if s.Server == nil {
		panic("服务配置信息不能为空")
	}
	if len(s.Server.ServerId) < 1 || len(s.Server.ServerType) < 1 {
		panic("服务器基本配置信息不能为空")
	}
	//初始化rpc服务
	s.Rpc = NewRpcServer(config.GetRpcConf(), s.Server)
	//订阅创建
	s.SubBuilder = NewRssBuilder(s.Server)
	//plugins
	for _, plugin := range s.plugins {
		plugin.Init(s)
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
	if err := s.Rpc.Subscribe(b.SetCodeType(CodeTypeProto).SetCallback(s.selfEventHandler).Build()); err != nil {
		panic(err)
	}
	//sub self json event
	b = s.SubBuilder.Build()
	if err := s.Rpc.Subscribe(b.SetSuffix(JsonSuffix).SetCodeType(CodeTypeJson).SetCallback(s.selfEventHandler).Build()); err != nil {
		panic(err)
	}
	//sub broadcast event
	b = s.SubBuilder.Build()
	if err := s.Rpc.SubscribeBroadcast(b.SetCodeType(CodeTypeProto).SetCallback(s.broadcastEventHandler).Build()); err != nil {
		panic(err)
	}
	//sub broadcast json event
	b = s.SubBuilder.Build()
	if err := s.Rpc.SubscribeBroadcast(b.SetSuffix(JsonSuffix).SetCodeType(CodeTypeJson).SetCallback(s.broadcastEventHandler).Build()); err != nil {
		panic(err)
	}

	//plugins
	for _, plugin := range s.plugins {
		plugin.AfterInit(s)
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
		plugin.BeforeShutdown(s)
	}
}

func (s *ServerBase) Shutdown() {
	//plugins
	for _, plugin := range s.plugins {
		plugin.Shutdown(s)
	}
	logger.Infof("shutdown service,type:%v,id:%v", s.Server.ServerType, s.Server.ServerId)
}

// HandleSelfEvent 内部事件处理
func (s *ServerBase) HandleSelfEvent(req *MsgRpc) []byte {
	resp, err := s.innerMsgHandler.DealMsg(CodeTypeProto, s.Rpc, req)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resp
}

func (s *ServerBase) HandleBroadcastEvent(req *MsgRpc) []byte {
	resp, err := s.innerMsgHandler.DealMsg(CodeTypeProto, s.Rpc, req)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resp
}
