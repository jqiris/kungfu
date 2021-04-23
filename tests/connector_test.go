package tests

import (
	"testing"

	"github.com/jqiris/kungfu/connector"
)

func TestConnector(t *testing.T) {
	cont := &connector.BaseConnector{}
	cont.Init()
}
