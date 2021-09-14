package stores

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/logger"
	"time"
)

var (
	defStoreKeeper StoreKeeper
)

func InitStoreKeeper(cfg config.StoresConf) {
	switch cfg.UseType {
	case "redis":
		defStoreKeeper = NewStoreRedis(
			WithRedisEndpoints(cfg.Endpoints),
			WithRedisDialTimeout(time.Duration(cfg.DialTimeout)*time.Second),
			WithRedisPassword(cfg.Password),
			WithRedisDB(cfg.DB),
			WithRedisPrefix(cfg.Prefix),
		)
	default:
		logger.Fatal("InitStoreKeeper failed")
	}
}

// StoreKeeper stores interface
type StoreKeeper interface {
	Set(key string, value interface{}, expire time.Duration) error
	SetNx(key string, value interface{}, expire time.Duration) error //set if not exist
	Get(key string, val interface{}) error
	GetInt(key string) int
	GetString(key string) string
	Del(keys ...string) error
	Exists(keys ...string) bool
	HSet(key, field string, val interface{}) error
	HGet(key, field string, val interface{}) error
	HGetAll(key string) (map[string]string, error)
	HDel(key string, fields ...string) error
	HDelAll(key string)
	HExists(key, field string) bool
	Expire(key string, expiration time.Duration) bool
	HKeys(key string) ([]string, error)
	LPush(key string, values ...interface{}) error
	RPush(key string, values ...interface{}) error
	LPop(key string, val interface{}) error
	RPop(key string, val interface{}) error
	BLPop(key string, val interface{}) error
	BRPop(key string, val interface{}) error
	LLen(key string) int64
	IsRedisNull(err error) bool
}

func Set(key string, value interface{}, expire time.Duration) error {
	return defStoreKeeper.Set(key, value, expire)
}
func SetNx(key string, value interface{}, expire time.Duration) error {
	return defStoreKeeper.SetNx(key, value, expire)
}
func Get(key string, val interface{}) error {
	return defStoreKeeper.Get(key, val)
}
func GetInt(key string) int {
	return defStoreKeeper.GetInt(key)
}
func GetString(key string) string {
	return defStoreKeeper.GetString(key)
}

func Del(keys ...string) error {
	return defStoreKeeper.Del(keys...)
}

func Exists(keys ...string) bool {
	return defStoreKeeper.Exists(keys...)
}

func HSet(key, field string, val interface{}) error {
	return defStoreKeeper.HSet(key, field, val)
}

func HGet(key, field string, val interface{}) error {
	return defStoreKeeper.HGet(key, field, val)
}

func HGetAll(key string) (map[string]string, error) {
	return defStoreKeeper.HGetAll(key)
}

func HDel(key string, fields ...string) error {
	return defStoreKeeper.HDel(key, fields...)
}

func HDelAll(key string) {
	defStoreKeeper.HDelAll(key)
}

func HExists(key, field string) bool {
	return defStoreKeeper.HExists(key, field)
}

func Expire(key string, expiration time.Duration) bool {
	return defStoreKeeper.Expire(key, expiration)
}

func HKeys(key string) ([]string, error) {
	return defStoreKeeper.HKeys(key)
}

func LPush(key string, values ...interface{}) error {
	return defStoreKeeper.LPush(key, values...)
}

func RPush(key string, values ...interface{}) error {
	return defStoreKeeper.RPush(key, values...)
}

func LPop(key string, val interface{}) error {
	return defStoreKeeper.LPop(key, val)
}

func RPop(key string, val interface{}) error {
	return defStoreKeeper.LPop(key, val)
}

func BLPop(key string, val interface{}) error {
	return defStoreKeeper.BLPop(key, val)
}

func BRPop(key string, val interface{}) error {
	return defStoreKeeper.BRPop(key, val)
}

func LLen(key string) int64 {
	return defStoreKeeper.LLen(key)
}

func IsRedisNull(err error) bool {
	return defStoreKeeper.IsRedisNull(err)
}
