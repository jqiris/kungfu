package stores

import (
	"context"
	"encoding/json"
	"github.com/jqiris/kungfu/logger"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type StoreRedis struct {
	Endpoints   []string
	Password    string
	DB          int
	DialTimeout time.Duration
	Client      *redis.Client
}

type StoreRedisOption func(s *StoreRedis)

func WithRedisEndpoints(endpoints []string) StoreRedisOption {
	return func(s *StoreRedis) {
		s.Endpoints = endpoints
	}
}

func WithRedisPassword(password string) StoreRedisOption {
	return func(s *StoreRedis) {
		s.Password = password
	}
}

func WithRedisDB(DB int) StoreRedisOption {
	return func(s *StoreRedis) {
		s.DB = DB
	}
}

func WithRedisDialTimeout(timeout time.Duration) StoreRedisOption {
	return func(s *StoreRedis) {
		s.DialTimeout = timeout
	}
}

func NewStoreRedis(opts ...StoreRedisOption) *StoreRedis {
	r := &StoreRedis{}
	for _, opt := range opts {
		opt(r)
	}
	c := redis.NewClient(&redis.Options{
		Addr:     strings.Join(r.Endpoints, ","),
		Password: r.Password,
		DB:       r.DB,
	})
	r.Client = c
	return r
}

func (s *StoreRedis) Set(key string, value interface{}, expire time.Duration) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.Set(ctx, key, bs, expire).Err()

}

func (s *StoreRedis) SetNx(key string, value interface{}, expire time.Duration) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.SetNX(ctx, key, bs, expire).Err()
}

func (s *StoreRedis) Get(key string, val interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	var bs []byte
	err := s.Client.Get(ctx, key).Scan(&bs)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, val)
}

func (s *StoreRedis) GetInt(key string) int {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	resp, err := s.Client.Get(ctx, key).Result()
	if err != nil {
		logger.Error(err)
		return 0
	}
	if res, err := strconv.Atoi(resp); err == nil {
		return res
	} else {
		logger.Error(err)
	}
	return 0
}

func (s *StoreRedis) GetString(key string) string {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	var bs []byte
	err := s.Client.Get(ctx, key).Scan(&bs)
	if err != nil {
		logger.Error(err)
		return ""
	}
	val := ""
	err = json.Unmarshal(bs, &val)
	if err != nil {
		logger.Error(err)
		return ""
	}
	return val
}

func (s *StoreRedis) Del(keys ...string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if _, err := s.Client.Del(ctx, keys...).Result(); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) Exists(keys ...string) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if val, err := s.Client.Exists(ctx, keys...).Result(); err != nil {
		logger.Error(err)
		return false
	} else {
		return val == 1
	}
}

func (s *StoreRedis) HSet(key, field string, val interface{}) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if _, err := s.Client.HSet(ctx, key, field, bs).Result(); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) HGet(key, field string, val interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	var bs []byte
	if err := s.Client.HGet(ctx, key, field).Scan(&bs); err != nil {
		return err
	}
	return json.Unmarshal(bs, val)
}

func (s *StoreRedis) HDel(key string, fields ...string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if _, err := s.Client.HDel(ctx, key, fields...).Result(); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) HDelAll(key string) {
	if fields, err := s.HKeys(key); err == nil {
		for _, field := range fields {
			if err = s.HDel(key, field); err != nil {
				logger.Error(err)
			}
		}
	} else {
		logger.Error(err)
	}
}

func (s *StoreRedis) HExists(key, field string) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if val, err := s.Client.HExists(ctx, key, field).Result(); err != nil {
		logger.Error(err)
		return false
	} else {
		return val
	}
}

func (s *StoreRedis) Expire(key string, expiration time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if val, err := s.Client.Expire(ctx, key, expiration).Result(); err != nil {
		logger.Error(err)
		return false
	} else {
		return val
	}
}

func (s *StoreRedis) HGetAll(key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.HGetAll(ctx, key).Result()
}

func (s *StoreRedis) HKeys(key string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.HKeys(ctx, key).Result()
}
