package tcpface

import "net"

// IConnection 定义连接接口
type IConnection interface {
	GetConn() net.Conn
	// GetConnID 获取当前连接ID
	GetConnID() int
	// RemoteAddr 获取远程客户端地址信息
	RemoteAddr() net.Addr
	Close() error
	StartWriter()
}

type IConnHandler func(server IServer, conn net.Conn, connId int) IConnection
