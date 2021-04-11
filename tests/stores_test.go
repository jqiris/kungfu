package tests

import (
	"github.com/jqiris/kungfu/stores"
	"testing"
	"time"
)

func TestStores(t *testing.T) {
	err := stores.Set("name", "jason", 5*time.Second)
	if err != nil {
		logger.Error(err)
		return
	}
	res, err := stores.Get("name")
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("Get name res:%+v", res)
	res2 := stores.GetString("name")
	logger.Infof("Get name res2:%+v", res2)

	res3 := stores.GetInt("name")
	logger.Infof("Get name res3:%+v", res3)
}
