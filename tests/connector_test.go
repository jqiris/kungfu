package tests

import (
	"testing"

	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/connector"
)

func TestConnector(t *testing.T) {
	confs := conf.GetConnectorConf()
	cont := &connector.BaseConnector{
		ConnectorConf: confs[0],
	}
	cont.Init()
}
