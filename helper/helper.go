package helper

import (
	"github.com/jqiris/kungfu/treaty"
)

func FindServerConfig(servers map[string]*treaty.Server, serverId string) *treaty.Server {
	if server, ok := servers[serverId]; ok {
		return server
	}
	return nil
}
