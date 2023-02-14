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

type IConnManager interface {
	Add(conn IConnection)                //添加链接
	Remove(conn IConnection)             //删除连接
	Get(connID int) (IConnection, error) //利用ConnID获取链接
	Len() int                            //获取当前连接
	ClearConn()                          //删除并停止所有链接
}
