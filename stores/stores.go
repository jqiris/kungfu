package stores

import (
	"time"

	"github.com/jqiris/kungfu/conf"
	"github.com/sirupsen/logrus"
)

var (
	logger         = logrus.WithField("package", "stores")
	defStoreKeeper StoreKeeper
)

func InitStoreKeeper(cfg conf.StoresConf) {
	switch cfg.UseType {
	case "redis":
		defStoreKeeper = NewStoreRedis(
			WithRedisEndpoints(cfg.Endpoints),
			WithRedisDialTimeout(time.Duration(cfg.DialTimeout)*time.Second),
			WithRedisPassword(cfg.Password),
			WithRedisDB(cfg.DB),
		)
	default:
		logger.Fatal("InitStoreKeeper failed")
	}
}

//stores interface
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
	HDel(key string, fields ...string) error
	HExists(key, field string) bool
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

func HDel(key string, fields ...string) error {
	return defStoreKeeper.HDel(key, fields...)
}

func HExists(key, field string) bool {
	return defStoreKeeper.HExists(key, field)
}
