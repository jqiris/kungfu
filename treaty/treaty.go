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
	Init()                     //初始化
	AfterInit()                //初始化后执行操作
	BeforeShutdown()           //服务关闭前操作
	Shutdown()                 //服务关闭操作
	GetServerId() int32        //获取server_id
	GetServerType() ServerType //获取服务类型
}
