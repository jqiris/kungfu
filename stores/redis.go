package stores

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
	"time"
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

func (s *StoreRedis) Get(key string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.Get(ctx, key).Result()
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
