package stores

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.Set(ctx, key, value, expire).Err()

}

func (s *StoreRedis) SetNx(key string, value interface{}, expire time.Duration) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.SetNX(ctx, key, value, expire).Err()
}

func (s *StoreRedis) Get(key string, val interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.Get(ctx, key).Scan(val)
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
	resp, err := s.Client.Get(ctx, key).Result()
	if err != nil {
		logger.Error(err)
		return ""
	}
	return resp
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

func (s *StoreRedis) HSet(key string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if _, err := s.Client.HSet(ctx, key, values...).Result(); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) HGet(key, field string, val interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if err := s.Client.HGet(ctx, key, field).Scan(val); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) HDel(key string, fields ...string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if _, err := s.Client.HDel(ctx, key, fields...).Result(); err != nil {
		return err
	}
	return nil
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
