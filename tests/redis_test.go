/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package tests

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/stores"
	"testing"
	"time"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func TestRedisHash(t *testing.T) {
	stores.HSet("a", "11", "111")
	stores.HSet("a", "22", "222")
	stores.HSet("a", "33", "333")
	res, err := stores.HGetAll("a")
	if err != nil {
		t.Fatal(err)
	}
	logger.Infof("res is: %+v", res)
	stores.HDel("a", "11")
	res, err = stores.HGetAll("a")
	if err != nil {
		t.Fatal(err)
	}
	logger.Infof("res is: %+v", res)

	//stores.HDelAll("a")
	if !stores.Expire("a", 2*time.Second) {
		t.Fatal("expire err")
	}
	logger.Info(stores.HExists("a", "11"))
	logger.Info(stores.HExists("a", "22"))
	select {
	case <-time.After(3 * time.Second):
		logger.Info(stores.Exists("a"))
		res, err = stores.HGetAll("a")
		if err != nil {
			t.Fatal(err)
		}
		logger.Infof("res is: %+v", res)
	}

}

func TestRedis(t *testing.T) {

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}

func TestRedis2(t *testing.T) {
	// SET key value EX 10 NX
	//set, err := rdb.SetNX(ctx, "key", "value", 10*time.Second).Result()
	//
	//// SET key value keepttl NX
	//set, err := rdb.SetNX(ctx, "key", "value", redis.KeepTTL).Result()
	//
	//// SORT list LIMIT 0 2 ASC
	//vals, err := rdb.Sort(ctx, "list", &redis.Sort{Offset: 0, Count: 2, Order: "ASC"}).Result()
	//
	//// ZRANGEBYSCORE zset -inf +inf WITHSCORES LIMIT 0 2
	//vals, err := rdb.ZRangeByScoreWithScores(ctx, "zset", &redis.ZRangeBy{
	//	Min: "-inf",
	//	Max: "+inf",
	//	Offset: 0,
	//	Count: 2,
	//}).Result()
	//
	//// ZINTERSTORE out 2 zset1 zset2 WEIGHTS 2 3 AGGREGATE SUM
	//vals, err := rdb.ZInterStore(ctx, "out", &redis.ZStore{
	//	Keys: []string{"zset1", "zset2"},
	//	Weights: []int64{2, 3}
	//}).Result()
	//
	//// EVAL "return {KEYS[1],ARGV[1]}" 1 "key" "hello"
	//vals, err := rdb.Eval(ctx, "return {KEYS[1],ARGV[1]}", []string{"key"}, "hello").Result()
	//
	//// custom command
	//res, err := rdb.Do(ctx, "set", "key", "value").Result()
}
