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
	JsonSuffix    = "json"
)
const (
	CodeTypeJson  = "json"
	CodeTypeProto = "proto"
)
