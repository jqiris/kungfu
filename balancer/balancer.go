package balancer

import (
	"github.com/jqiris/kungfu/treaty"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.WithField("package", "balancer")
)

type Balancer interface {
	treaty.ServerEntity
	Balance(remoteAddr string) (*treaty.Server, error)
}
