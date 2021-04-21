package conf

import "github.com/jqiris/zinx/utils"

type Config struct {
	Discover  DiscoverConf    `yaml:"discover"`
	Rpcx      RpcxConf        `yaml:"rpcx"`
	Stores    StoresConf      `yaml:"stores"`
	Coder     CoderConf       `yaml:"coder"`
	Connector utils.GlobalObj `yaml:"connector"`
}

type DiscoverConf struct {
	UseType     string   `yaml:"use_type"`
	DialTimeout int      `yaml:"dial_timeout"`
	Endpoints   []string `yaml:"endpoints"`
}

type RpcxConf struct {
	UseType     string   `yaml:"use_type"`
	DialTimeout int      `yaml:"dial_timeout"`
	Endpoints   []string `yaml:"endpoints"`
}

type StoresConf struct {
	UseType     string   `yaml:"use_type"`
	DialTimeout int      `yaml:"dial_timeout"`
	Endpoints   []string `yaml:"endpoints"`
	Password    string   `yaml:"password"`
	DB          int      `yaml:"db"`
}

type CoderConf struct {
	UseType string `yaml:"use_type"`
}

//type ConnectorConf struct {
//	Host             string `yaml:"host"`                //当前服务器主机IP
//	TCPPort          int    `yaml:"tcp_port"`            //当前服务器主机监听端口号
//	Name             string `yaml:"name"`                //当前服务器名称
//	MaxPacketSize    uint32 `yaml:"max_packet_size"`     //都需数据包的最大值
//	MaxConn          int    `yaml:"max_conn"`            //当前服务器主机允许的最大链接个数
//	WorkerPoolSize   uint32 `yaml:"worker_pool_size"`    //业务工作Worker池的数量
//	MaxWorkerTaskLen uint32 `yaml:"max_worker_task_len"` //业务工作Worker对应负责的任务队列最大任务存储数量
//	MaxMsgChanLen    uint32 `yaml:"max_msg_chan_len"`    //SendBuffMsg发送消息的缓冲最大长度
//	LogDir           string `yaml:"log_dir"`             //日志所在文件夹 默认"./log"
//	LogFile          string `yaml:"log_file"`            //日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr
//	LogDebugClose    bool   `yaml:"log_debug_close"`     //是否关闭Debug日志级别调试信息 默认false  -- 默认打开debug信息
//}
