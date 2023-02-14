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

package rpc

import "github.com/jqiris/kungfu/v2/treaty"

type ServerCreator func(s *treaty.Server) (ServerEntity, error)

// ServerEntity server entity
type ServerEntity interface {
	Init()                                   //初始化
	AfterInit()                              //初始化后执行操作
	BeforeShutdown()                         //服务关闭前操作
	Shutdown()                               //服务关闭操作
	HandleSelfEvent(req *MsgRpc) []byte      //处理自己的事件
	HandleBroadcastEvent(req *MsgRpc) []byte //处理广播事件
}

//ServerPlugin server extand
type ServerPlugin interface {
	Init(s *ServerBase)           //初始化
	AfterInit(s *ServerBase)      //初始化后执行操作
	BeforeShutdown(s *ServerBase) //服务关闭前操作
	Shutdown(s *ServerBase)       //服务关闭操作
}
