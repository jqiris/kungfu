package treaty

import (
	"encoding/json"
	"fmt"
)

//const
const (
	MinServerId = 1000
)

const (
	DiscoverEtcd = "etcd"
)

//server type name
type ServerType int

const (
	ServerNone      ServerType = iota + 1000
	ServerGate                 // gate服务器
	ServerConnector            // connector服务器
	ServerGame                 //游戏服务器
)

//server struct
type Server struct {
	ServerId   int        `json:"server_id"`   //服务器ID
	ServerType ServerType `json:"server_type"` //服务器类型
	ServerName string     `json:"server_name"` //服务器名字
	ServerHost string     `json:"server_host"` //服务器地址
}

func (s *Server) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		logger.Error(err)
		return ""
	}
	return string(b)
}

func (s *Server) RegId() string {
	return fmt.Sprintf("%d/%d", s.ServerType, s.ServerId)
}

func (s *Server) RegType() string {
	return fmt.Sprintf("%d", s.ServerType)
}
