package tcpface

// IMsgHandle 消息管理抽象层
type IMsgHandle interface {
	StartWorkerPool() //启动worker工作池
	Handle(iConn IConnection)
}
