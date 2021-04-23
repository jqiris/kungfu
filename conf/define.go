package conf

import "github.com/jqiris/zinx/utils"

type Config struct {
	Discover  DiscoverConf      `yaml:"discover"`
	Rpcx      RpcxConf          `yaml:"rpcx"`
	Stores    StoresConf        `yaml:"stores"`
	Coder     CoderConf         `yaml:"coder"`
	Balancer  []BalancerConf    `yaml:"balancer"`
	Connector []utils.GlobalObj `yaml:"connector"`
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

type BalancerConf struct {
	ServerId int    `yaml:"server_id"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	HttpPort int    `yaml:"http_port"`
}
