package tests

import (
	"fmt"
	"github.com/jqiris/kungfu/logger"
	"testing"
	"time"

	"github.com/jqiris/kungfu/stores"
)

func TestStores(t *testing.T) {
	err := stores.Set("name", "jason", 5*time.Second)
	if err != nil {
		logger.Error(err)
		return
	}
	var res string
	err = stores.Get("name", &res)
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

func TestStoreList(t *testing.T) {
	key := "myList"
	var err error
	err = stores.LPush(key, 1, 2, 3, 4, 5)
	if err != nil {
		logger.Error(err)
		return
	}
	fmt.Println("length:", stores.LLen(key))
	var a string
	if a, err = stores.BRPopString(key); err != nil {
		logger.Error(err)
		return
	}
	fmt.Println("pop:", a)
	fmt.Println("length:", stores.LLen(key))

	if err = stores.BRPop(key, &a); err != nil {
		logger.Error(err)
		return
	}
	fmt.Println("pop:", a)
	fmt.Println("length:", stores.LLen(key))
	select {}
}
