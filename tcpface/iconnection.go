/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

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
