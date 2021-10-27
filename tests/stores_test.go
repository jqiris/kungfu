package tests

import (
	"fmt"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/utils"
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

func TestZAdd(t *testing.T) {
	//rand.Seed(time.Now().UnixNano())
	//maxNum := 200000
	//list := make([]*redis.Z, 0)
	//for i := 0; i < maxNum; i++ {
	//	//score, member := float64(rand.Intn(maxNum)), fmt.Sprintf("member_%d", i)
	//	score, member := float64(rand.Intn(maxNum)), i+1
	//	list = append(list, &redis.Z{
	//		Score:  score,
	//		Member: member,
	//	})
	//}
	//if err := stores.ZAdd("list_rank", list...); err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("ZAdd成功")
	//list, err := stores.ZRevRangeWithScores("list_rank", 0, 100)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(list)
	//fmt.Println(len(list))
	if err := stores.ZRem("list_rank", utils.IntToString(6113)); err != nil {
		fmt.Println("111", err)
	}
	if rank, err := stores.ZRevRank("list_rank", utils.IntToString(6113)); err != nil {
		fmt.Println("222", err)
	} else {
		fmt.Println(rank)
	}
	//if s, err := stores.ZIncrBy("list_rank", 3, utils.IntToString(6113)); err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(s)
	//}
	//score := stores.ZScore("list_rank", utils.IntToString(6113))
	//fmt.Println(score)
}
