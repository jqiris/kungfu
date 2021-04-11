package stores

import (
	"github.com/jqiris/kungfu/conf"
	"github.com/sirupsen/logrus"
	"time"
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
	Get(key string) (interface{}, error)
	GetInt(key string) int
	GetString(key string) string
}

func Set(key string, value interface{}, expire time.Duration) error {
	return defStoreKeeper.Set(key, value, expire)
}
func SetNx(key string, value interface{}, expire time.Duration) error {
	return defStoreKeeper.SetNx(key, value, expire)
}
func Get(key string) (interface{}, error) {
	return defStoreKeeper.Get(key)
}
func GetInt(key string) int {
	return defStoreKeeper.GetInt(key)
}
func GetString(key string) string {
	return defStoreKeeper.GetString(key)
}
