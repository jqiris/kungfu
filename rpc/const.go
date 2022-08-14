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
	DefaultExName = ""
	DefaultExType = "direct"
	FanoutExType  = "fanout"
	DefaultRtKey  = ""
	JsonSuffix    = "json"
)
const (
	CodeTypeJson  = "json"
	CodeTypeProto = "proto"
)
