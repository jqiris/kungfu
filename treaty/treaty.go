package treaty

import (
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.WithField("package", "treaty")
)

const (
	MinServerId       = 1000
	MessageHeaderSize = 20
)

//server entity
type ServerEntity interface {
	Init()                                                    //初始化
	AfterInit()                                               //初始化后执行操作
	BeforeShutdown()                                          //服务关闭前操作
	Shutdown()                                                //服务关闭操作
	GetServer() *Server                                       //获取服务
	RegEventHandlerSelf(handler func(req []byte) []byte)      //注册自己事件处理器
	RegEventHandlerBroadcast(handler func(req []byte) []byte) //注册广播事件处理器
	SetServerId(serverId int32)                               //设置SetServerId
	GetServerId() int32                                       //获取serverId
}
