package utils

import (
	"github.com/jqiris/kungfu/logger"
	"strconv"
)

func StringToInt(s string) int {
	if res, err := strconv.Atoi(s); err != nil {
		logger.Error(err)
		return 0
	} else {
		return res
	}
}
