package conf

import (
	"github.com/jqiris/kungfu/treaty"
)

type Config struct {
	Discover  DiscoverConf              `json:"discover"`
	Rpcx      RpcxConf                  `json:"rpcx"`
	Stores    StoresConf                `json:"stores"`
	Coder     CoderConf                 `json:"coder"`
	Connector ConnectorConf             `json:"connector"`
	Servers   map[string]*treaty.Server `json:"servers"`
	Launch    []string                  `json:"launch"`
}

type DiscoverConf struct {
	UseType     string   `json:"use_type"`
	DialTimeout int      `json:"dial_timeout"`
	Endpoints   []string `json:"endpoints"`
}

type RpcxConf struct {
	UseType     string   `json:"use_type"`
	DialTimeout int      `json:"dial_timeout"`
	Endpoints   []string `json:"endpoints"`
}

type StoresConf struct {
	UseType     string   `json:"use_type"`
	DialTimeout int      `json:"dial_timeout"`
	Endpoints   []string `json:"endpoints"`
	Password    string   `json:"password"`
	DB          int      `json:"db"`
}

type CoderConf struct {
	UseType string `json:"use_type"`
}

type ConnectorConf struct {
	Version          string `json:"version"`             //当前Zinx版本号
	MaxPacketSize    int32  `json:"max_packet_size"`     //都需数据包的最大值
	MaxConn          int    `json:"max_conn"`            //当前服务器主机允许的最大链接个数
	WorkerPoolSize   int32  `json:"worker_pool_size"`    //业务工作Worker池的数量
	MaxWorkerTaskLen int32  `json:"max_worker_task_len"` //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    int32  `json:"max_msg_chan_len"`    //SendBuffMsg发送消息的缓冲最大长度
	LogDir           string `json:"log_dir"`             //日志所在文件夹 默认"./log"
	LogFile          string `json:"log_file"`            //日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr
	LogDebugClose    bool   `json:"log_debug_close"`     //是否关闭Debug日志级别调试信息 默认false  -- 默认打开debug信息
	TokenKey         string `json:"token_key"`           //token生成key
}
