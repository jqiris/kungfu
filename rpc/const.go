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

const (
	Balancer  = "balancer"
	Connector = "connector"
	Server    = "backend"
	Database  = "database"
	Loader    = "loader"
)
const (
	DefaultQueue  = "dq"
	DefaultSuffix = ""
	DefaultExName = "exchange"
	DefaultExType = "direct"
	FanoutExType  = "fanout"
	TopicEXType   = "topic"
	DefaultRtKey  = ""
	DefaultReply  = "reply"
	JsonSuffix    = "json"
)
const (
	CodeTypeJson  = "json"
	CodeTypeProto = "proto"
)
