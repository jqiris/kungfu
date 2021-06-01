package utils

import (
	"github.com/sirupsen/logrus"
	"strconv"
)

var (
	logger = logrus.WithField("package", "utils")
)

func StringToInt(s string) int {
	if res, err := strconv.Atoi(s); err != nil {
		logger.Error(err)
		return 0
	} else {
		return res
	}
}
