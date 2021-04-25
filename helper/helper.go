package helper

import (
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/zinx/utils"
)

func FindServerConfig(servers []*treaty.Server, serverId int32) *treaty.Server {
	for _, server := range servers {
		if server.ServerId == serverId {
			return server
		}
	}
	return nil
}

func FindConnectorConfig(servers []*utils.GlobalObj, serverId int32) *utils.GlobalObj {
	for _, server := range servers {
		if server.ServerId == serverId {
			return server
		}
	}
	return nil
}
