package stores

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jqiris/kungfu/logger"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	defPrefix = "store"
)

type StoreRedis struct {
	Endpoints   []string
	Password    string
	DB          int
	DialTimeout time.Duration
	Client      *redis.Client
	Prefix      string
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
func WithRedisPrefix(prefix string) StoreRedisOption {
	return func(s *StoreRedis) {
		s.Prefix = prefix
	}
}

func NewStoreRedis(opts ...StoreRedisOption) *StoreRedis {
	r := &StoreRedis{
		Prefix: defPrefix,
	}
	for _, opt := range opts {
		opt(r)
	}
	if len(r.Prefix) == 0 {
		r.Prefix = defPrefix
	}
	c := redis.NewClient(&redis.Options{
		Addr:     strings.Join(r.Endpoints, ","),
		Password: r.Password,
		DB:       r.DB,
	})
	r.Client = c
	return r
}

func (s *StoreRedis) GetKey(key string) string {
	return s.Prefix + ":" + key
}

func (s *StoreRedis) GetKeys(keys []string) []string {
	var list []string
	for _, v := range keys {
		list = append(list, s.GetKey(v))
	}
	return list
}

func (s *StoreRedis) GetValues(values []interface{}) []interface{} {
	res := make([]interface{}, 0)
	for _, v := range values {
		if bs, err := json.Marshal(v); err == nil {
			res = append(res, bs)
		}
	}
	return res
}

func (s *StoreRedis) Set(key string, value interface{}, expire time.Duration) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.Set(ctx, s.GetKey(key), bs, expire).Err()

}

func (s *StoreRedis) SetNx(key string, value interface{}, expire time.Duration) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.SetNX(ctx, s.GetKey(key), bs, expire).Err()
}

func (s *StoreRedis) Get(key string, val interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	var bs []byte
	err := s.Client.Get(ctx, s.GetKey(key)).Scan(&bs)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, val)
}

func (s *StoreRedis) GetInt(key string) int {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	resp, err := s.Client.Get(ctx, s.GetKey(key)).Result()
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
	err := s.Client.Get(ctx, s.GetKey(key)).Scan(&bs)
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
	if _, err := s.Client.Del(ctx, s.GetKeys(keys)...).Result(); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) Exists(keys ...string) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if val, err := s.Client.Exists(ctx, s.GetKeys(keys)...).Result(); err != nil {
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
	if _, err := s.Client.HSet(ctx, s.GetKey(key), field, bs).Result(); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) HGet(key, field string, val interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	var bs []byte
	if err := s.Client.HGet(ctx, s.GetKey(key), field).Scan(&bs); err != nil {
		return err
	}
	return json.Unmarshal(bs, val)
}

func (s *StoreRedis) HDel(key string, fields ...string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if _, err := s.Client.HDel(ctx, s.GetKey(key), fields...).Result(); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) HDelAll(key string) {
	nKey := s.GetKey(key)
	if fields, err := s.HKeys(nKey); err == nil {
		for _, field := range fields {
			if err = s.HDel(nKey, field); err != nil {
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
	if val, err := s.Client.HExists(ctx, s.GetKey(key), field).Result(); err != nil {
		logger.Error(err)
		return false
	} else {
		return val
	}
}

func (s *StoreRedis) Expire(key string, expiration time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if val, err := s.Client.Expire(ctx, s.GetKey(key), expiration).Result(); err != nil {
		logger.Error(err)
		return false
	} else {
		return val
	}
}

func (s *StoreRedis) HGetAll(key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.HGetAll(ctx, s.GetKey(key)).Result()
}

func (s *StoreRedis) HKeys(key string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.HKeys(ctx, s.GetKey(key)).Result()
}

func (s *StoreRedis) LPush(key string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if _, err := s.Client.LPush(ctx, s.GetKey(key), s.GetValues(values)...).Result(); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) RPush(key string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if _, err := s.Client.RPush(ctx, s.GetKey(key), s.GetValues(values)...).Result(); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) LPop(key string, val interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	var bs []byte
	err := s.Client.LPop(ctx, s.GetKey(key)).Scan(&bs)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, val)
}
func (s *StoreRedis) RPop(key string, val interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	var bs []byte
	err := s.Client.RPop(ctx, s.GetKey(key)).Scan(&bs)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, val)
}

func (s *StoreRedis) BLPop(key string, val interface{}) error {
	bss, err := s.Client.BLPop(context.Background(), s.DialTimeout, s.GetKey(key)).Result()
	if err != nil {
		return err
	}
	if len(bss) < 2 {
		return errors.New("wait timeout")
	}
	return json.Unmarshal([]byte(bss[1]), val)
}

func (s *StoreRedis) BRPop(key string, val interface{}) error {
	bss, err := s.Client.BRPop(context.Background(), s.DialTimeout, s.GetKey(key)).Result()
	if err != nil {
		return err
	}
	if len(bss) < 2 {
		return errors.New("wait timeout")
	}
	return json.Unmarshal([]byte(bss[1]), val)
}

func (s *StoreRedis) BLPopString(key string) (string, error) {
	bss, err := s.Client.BLPop(context.Background(), s.DialTimeout, s.GetKey(key)).Result()
	if err != nil {
		return "", err
	}
	if len(bss) < 2 {
		return "", errors.New("wait timeout")
	}
	return bss[1], nil
}

func (s *StoreRedis) BRPopString(key string) (string, error) {
	bss, err := s.Client.BRPop(context.Background(), s.DialTimeout, s.GetKey(key)).Result()
	if err != nil {
		return "", err
	}
	if len(bss) < 2 {
		return "", errors.New("wait timeout")
	}
	return bss[1], nil
}

func (s *StoreRedis) LLen(key string) int64 {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.LLen(ctx, s.GetKey(key)).Val()
}

func (s *StoreRedis) IsRedisNull(err error) bool {
	return err == redis.Nil
}
