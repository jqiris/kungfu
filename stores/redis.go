package stores

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jqiris/kungfu/v2/logger"

	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/golang/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
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
	lock        *sync.Mutex
	rs          *redsync.Redsync
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
	pool := goredis.NewPool(c)
	r.rs = redsync.New(pool)
	r.lock = new(sync.Mutex)
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

func (s *StoreRedis) GetValues(values []any) []any {
	res := make([]any, 0)
	for _, v := range values {
		if bs, err := json.Marshal(v); err == nil {
			res = append(res, bs)
		}
	}
	return res
}

func (s *StoreRedis) Set(key string, value any, expire time.Duration) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.Set(ctx, s.GetKey(key), bs, expire).Err()
}
func (s *StoreRedis) SetProto(key string, value proto.Message, expire time.Duration) error {
	bs, err := proto.Marshal(value)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.Set(ctx, s.GetKey(key), bs, expire).Err()
}

func (s *StoreRedis) SetNx(key string, value any, expire time.Duration) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.SetNX(ctx, s.GetKey(key), bs, expire).Err()
}
func (s *StoreRedis) SetProtoNx(key string, value proto.Message, expire time.Duration) error {
	bs, err := proto.Marshal(value)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.SetNX(ctx, s.GetKey(key), bs, expire).Err()
}

func (s *StoreRedis) Get(key string, val any) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	bs, err := s.Client.Get(ctx, s.GetKey(key)).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, val)
}

func (s *StoreRedis) GetProto(key string, val proto.Message) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	bs, err := s.Client.Get(ctx, s.GetKey(key)).Bytes()
	if err != nil {
		return err
	}
	return proto.Unmarshal(bs, val)
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
	bs, err := s.Client.Get(ctx, s.GetKey(key)).Bytes()
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

func (s *StoreRedis) Del(keys ...string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.Del(ctx, s.GetKeys(keys)...).Result()
}

func (s *StoreRedis) DelPattern(pattern string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	keys := s.Client.Keys(ctx, s.GetKey(pattern)).Val()
	if len(keys) > 0 {
		return s.Client.Del(ctx, keys...).Result()
	}
	return 0, nil
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

func (s *StoreRedis) HSet(key, field string, val any) error {
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

func (s *StoreRedis) HSetNx(key, field string, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if _, err := s.Client.HSetNX(ctx, s.GetKey(key), field, bs).Result(); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) HIncrBy(key, field string, incr int64) int64 {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.HIncrBy(ctx, s.GetKey(key), field, incr).Val()
}

func (s *StoreRedis) HGet(key, field string, val any) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	var bs []byte
	var err error
	if bs, err = s.Client.HGet(ctx, s.GetKey(key), field).Bytes(); err != nil {
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

func (s *StoreRedis) LPush(key string, values ...any) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if _, err := s.Client.LPush(ctx, s.GetKey(key), s.GetValues(values)...).Result(); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) RPush(key string, values ...any) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	if _, err := s.Client.RPush(ctx, s.GetKey(key), s.GetValues(values)...).Result(); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) LPop(key string, val any) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	bs, err := s.Client.LPop(ctx, s.GetKey(key)).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, val)
}
func (s *StoreRedis) RPop(key string, val any) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	bs, err := s.Client.RPop(ctx, s.GetKey(key)).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, val)
}

func (s *StoreRedis) BLPop(key string, val any) error {
	bss, err := s.Client.BLPop(context.Background(), s.DialTimeout, s.GetKey(key)).Result()
	if err != nil {
		return err
	}
	if len(bss) < 2 {
		return errors.New("wait timeout")
	}
	return json.Unmarshal([]byte(bss[1]), val)
}

func (s *StoreRedis) BRPop(key string, val any) error {
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
func (s *StoreRedis) FlushDB() error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.FlushDB(ctx).Err()
}

func (s *StoreRedis) FlushDBAsync() error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.FlushDBAsync(ctx).Err()
}

func (s *StoreRedis) FlushAll() error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.FlushAll(ctx).Err()
}

func (s *StoreRedis) FlushAllAsync() error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.FlushAllAsync(ctx).Err()
}

func (s *StoreRedis) SAdd(key string, members ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.SAdd(ctx, s.GetKey(key), members...).Err()
}

func (s *StoreRedis) SCard(key string) int64 {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.SCard(ctx, s.GetKey(key)).Val()
}

func (s *StoreRedis) SRem(key string, members ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.SRem(ctx, s.GetKey(key), members...).Err()
}

func (s *StoreRedis) SMembers(key string) []string {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.SMembers(ctx, s.GetKey(key)).Val()
}

func (s *StoreRedis) SRandMember(key string) string {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.SRandMember(ctx, s.GetKey(key)).Val()
}

func (s *StoreRedis) SRandMemberN(key string, count int64) []string {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.SRandMemberN(ctx, s.GetKey(key), count).Val()
}

func (s *StoreRedis) SIsMember(key string, member interface{}) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.SIsMember(ctx, s.GetKey(key), member).Val()
}

func (s *StoreRedis) ZAdd(key string, members ...*redis.Z) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.ZAdd(ctx, s.GetKey(key), members...).Err()
}

func (s *StoreRedis) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.ZRangeWithScores(ctx, s.GetKey(key), start, stop).Result()
}

func (s *StoreRedis) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.ZRevRangeWithScores(ctx, s.GetKey(key), start, stop).Result()
}

func (s *StoreRedis) ZRevRank(key, member string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.ZRevRank(ctx, s.GetKey(key), member).Result()
}

func (s *StoreRedis) ZScore(key, member string) float64 {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.ZScore(ctx, s.GetKey(key), member).Val()
}

func (s *StoreRedis) ZIncrBy(key string, increment float64, member string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.ZIncrBy(ctx, s.GetKey(key), increment, member).Result()
}

func (s *StoreRedis) ZRem(key string, members ...any) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.ZRem(ctx, s.GetKey(key), members...).Err()
}

func (s *StoreRedis) ZCard(key string) int64 {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.ZCard(ctx, s.GetKey(key)).Val()
}

func (s *StoreRedis) Incr(key string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.Incr(ctx, s.GetKey(key)).Result()
}

func (s *StoreRedis) Decr(key string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.DialTimeout)
	defer cancel()
	return s.Client.Decr(ctx, s.GetKey(key)).Result()
}

func (s *StoreRedis) Lock(key string) (*redsync.Mutex, context.Context, error) {
	mutex := s.rs.NewMutex(s.GetKey(key))
	ctx := context.Background()
	if err := mutex.LockContext(ctx); err != nil {
		return mutex, ctx, err
	}
	return mutex, ctx, nil
}

func (s *StoreRedis) Unlock(mutex *redsync.Mutex, ctx context.Context) error {
	if _, err := mutex.UnlockContext(ctx); err != nil {
		return err
	}
	return nil
}

func (s *StoreRedis) TxPipeline() redis.Pipeliner {
	return s.Client.TxPipeline()
}
