package config

import (
	"github.com/jqiris/kungfu/treaty"
)

type Config struct {
	Discover  DiscoverConf              `json:"discover"`
	RpcX      RpcXConf                  `json:"rpcx"`
	Stores    StoresConf                `json:"stores"`
	Connector ConnectorConf             `json:"connector"`
	Servers   map[string]*treaty.Server `json:"servers"`
	Launch    []string                  `json:"launch"`
}

type DiscoverConf struct {
	UseType     string   `json:"use_type"`
	DialTimeout int      `json:"dial_timeout"`
	Endpoints   []string `json:"endpoints"`
}

type RpcXConf struct {
	UseType     string   `json:"use_type"`
	DialTimeout int      `json:"dial_timeout"`
	Endpoints   []string `json:"endpoints"`
	DebugMsg    bool     `json:"debug_msg"`
}

type StoresConf struct {
	UseType     string   `json:"use_type"`
	DialTimeout int      `json:"dial_timeout"`
	Endpoints   []string `json:"endpoints"`
	Password    string   `json:"password"`
	DB          int      `json:"db"`
}

type ConnectorConf struct {
	UseType           string `json:"use_type"`            //使用的协议
	UseWebsocket      bool   `json:"use_websocket"`       //是否使用websocket
	WebsocketPath     string `json:"websocket_path"`      //websocket路径
	UseSerializer     string `json:"use_serializer"`      //使用的协议
	ProtoPath         string `json:"proto_path"`          //protobuf位置
	HeartbeatInterval int    `json:"heartbeat_interval"`  //心跳间隔
	Version           string `json:"version"`             //当前tcpserver版本号
	MaxPacketSize     int32  `json:"max_packet_size"`     //都需数据包的最大值
	MaxConn           int    `json:"max_conn"`            //当前服务器主机允许的最大链接个数
	WorkerPoolSize    int    `json:"worker_pool_size"`    //业务工作Worker池的数量
	MaxWorkerTaskLen  int32  `json:"max_worker_task_len"` //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen     int32  `json:"max_msg_chan_len"`    //SendBuffMsg发送消息的缓冲最大长度
	LogDir            string `json:"log_dir"`             //日志所在文件夹 默认"./log"
	LogFile           string `json:"log_file"`            //日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr
	LogDebugClose     bool   `json:"log_debug_close"`     //是否关闭Debug日志级别调试信息 默认false  -- 默认打开debug信息
	TokenKey          string `json:"token_key"`           //token生成key
}
