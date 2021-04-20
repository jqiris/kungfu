package connector

import (
	"github.com/jqiris/kungfu/treaty"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.WithField("package", "connector")
)

type Connector interface {
	treaty.ServerEntity
}
