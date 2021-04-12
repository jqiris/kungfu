package treaty

import "github.com/sirupsen/logrus"

var (
	logger = logrus.WithField("package", "treaty")
)

//server entity
type ServerEntity interface {
	OnInit() error
	OnRun() error
	OnRegister() error
	UnRegister() error
	OnStop() error
}
