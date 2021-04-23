package backend

import (
	"github.com/jqiris/kungfu/treaty"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.WithField("package", "backend")
)

type BackEnd interface {
	treaty.ServerEntity
}
