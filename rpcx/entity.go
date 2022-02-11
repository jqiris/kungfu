package rpcx

import "github.com/jqiris/kungfu/treaty"

type ServerCreator func(s *treaty.Server) (ServerEntity, error)

// ServerEntity server entity
type ServerEntity interface {
	Init()           //初始化
	AfterInit()      //初始化后执行操作
	BeforeShutdown() //服务关闭前操作
	Shutdown()       //服务关闭操作
}
