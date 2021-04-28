package conf

import (
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/zinx/utils"
)

type Config struct {
	Discover  DiscoverConf       `json:"discover"`
	Rpcx      RpcxConf           `json:"rpcx"`
	Stores    StoresConf         `json:"stores"`
	Coder     CoderConf          `json:"coder"`
	Balancer  []*treaty.Server   `json:"balancer"`
	Connector []*utils.GlobalObj `json:"connector"`
	Servers   []*treaty.Server   `json:"backend"`
	Launch    []int32            `json:"launch"`
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
