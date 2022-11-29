package stores

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/golang/protobuf/proto"
	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/logger"
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
	Set(key string, value any, expire time.Duration) error
	SetNx(key string, value any, expire time.Duration) error //set if not exist
	SetProto(key string, value proto.Message, expire time.Duration) error
	SetProtoNx(key string, value proto.Message, expire time.Duration) error
	Get(key string, val any) error
	GetInt(key string) int
	GetString(key string) string
	GetProto(key string, val proto.Message) error
	Del(keys ...string) (int64, error)
	DelPattern(pattern string) (int64, error)
	Exists(keys ...string) bool
	HSet(key, field string, val any) error
	HSetNx(key, field string, val any) error
	HIncrBy(key, field string, incr int64) int64
	HGet(key, field string, val any) error
	HGetAll(key string) (map[string]string, error)
	HDel(key string, fields ...string) error
	HDelAll(key string)
	HExists(key, field string) bool
	Expire(key string, expiration time.Duration) bool
	HKeys(key string) ([]string, error)
	SAdd(key string, members ...interface{}) error
	SCard(key string) int64
	SRem(key string, members ...interface{}) error
	SMembers(key string) []string
	SRandMember(key string) string
	SRandMemberN(key string, count int64) []string
	SPop(key string) string
	SPopN(key string, count int64) []string
	SIsMember(key string, member interface{}) bool
	ZAdd(key string, members ...*redis.Z) error
	ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error)
	ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error)
	ZRank(key, member string) (int64, error)
	ZRevRank(key, member string) (int64, error)
	ZScore(key, member string) float64
	ZIncrBy(key string, increment float64, member string) (float64, error)
	ZRem(key string, members ...any) error
	ZRemRangeByScore(key string, min, max string) error
	ZCard(key string) int64
	LPush(key string, values ...any) error
	RPush(key string, values ...any) error
	LPop(key string, val any) error
	RPop(key string, val any) error
	BLPop(key string, val any) error
	BRPop(key string, val any) error
	BLPopString(key string) (string, error)
	BRPopString(key string) (string, error)
	Incr(key string) (int64, error)
	Decr(key string) (int64, error)
	LLen(key string) int64
	IsRedisNull(err error) bool
	FlushDB() error
	FlushDBAsync() error
	FlushAll() error
	FlushAllAsync() error
	Lock(key string) (*redsync.Mutex, context.Context, error)
	Unlock(mutex *redsync.Mutex, ctx context.Context) error
	TxPipeline() redis.Pipeliner
	GetKey(key string) string
}

func GetDefStoreKeeper() StoreKeeper {
	return defStoreKeeper
}

func Set(key string, value any, expire time.Duration) error {
	return defStoreKeeper.Set(key, value, expire)
}
func SetNx(key string, value any, expire time.Duration) error {
	return defStoreKeeper.SetNx(key, value, expire)
}
func SetProto(key string, value proto.Message, expire time.Duration) error {
	return defStoreKeeper.SetProto(key, value, expire)
}
func SetProtoNx(key string, value proto.Message, expire time.Duration) error {
	return defStoreKeeper.SetProtoNx(key, value, expire)
}
func Get(key string, val any) error {
	return defStoreKeeper.Get(key, val)
}
func GetInt(key string) int {
	return defStoreKeeper.GetInt(key)
}
func GetString(key string) string {
	return defStoreKeeper.GetString(key)
}
func GetProto(key string, val proto.Message) error {
	return defStoreKeeper.GetProto(key, val)
}

func Del(keys ...string) (int64, error) {
	return defStoreKeeper.Del(keys...)
}

func DelPattern(pattern string) (int64, error) {
	return defStoreKeeper.DelPattern(pattern)
}

func Exists(keys ...string) bool {
	return defStoreKeeper.Exists(keys...)
}

func HSet(key, field string, val any) error {
	return defStoreKeeper.HSet(key, field, val)
}
func HSetNx(key, field string, val any) error {
	return defStoreKeeper.HSetNx(key, field, val)
}

func HIncrBy(key, field string, incr int64) int64 {
	return defStoreKeeper.HIncrBy(key, field, incr)
}

func HGet(key, field string, val any) error {
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

func LPush(key string, values ...any) error {
	return defStoreKeeper.LPush(key, values...)
}

func RPush(key string, values ...any) error {
	return defStoreKeeper.RPush(key, values...)
}

func LPop(key string, val any) error {
	return defStoreKeeper.LPop(key, val)
}

func RPop(key string, val any) error {
	return defStoreKeeper.LPop(key, val)
}

func BLPop(key string, val any) error {
	return defStoreKeeper.BLPop(key, val)
}

func BRPop(key string, val any) error {
	return defStoreKeeper.BRPop(key, val)
}

func BLPopString(key string) (string, error) {
	return defStoreKeeper.BLPopString(key)
}

func BRPopString(key string) (string, error) {
	return defStoreKeeper.BRPopString(key)
}

func LLen(key string) int64 {
	return defStoreKeeper.LLen(key)
}

func IsRedisNull(err error) bool {
	return defStoreKeeper.IsRedisNull(err)
}
func FlushDBAsync() error {
	return defStoreKeeper.FlushDBAsync()
}
func FlushDB() error {
	return defStoreKeeper.FlushDB()
}

func FlushAllAsync() error {
	return defStoreKeeper.FlushAllAsync()
}
func FlushAll() error {
	return defStoreKeeper.FlushAll()
}

func SAdd(key string, members ...interface{}) error {
	return defStoreKeeper.SAdd(key, members...)
}
func SCard(key string) int64 {
	return defStoreKeeper.SCard(key)
}

func SRem(key string, members ...interface{}) error {
	return defStoreKeeper.SRem(key, members...)
}

func SMembers(key string) []string {
	return defStoreKeeper.SMembers(key)
}

func SRandMember(key string) string {
	return defStoreKeeper.SRandMember(key)
}

func SRandMemberN(key string, count int64) []string {
	return defStoreKeeper.SRandMemberN(key, count)
}

func SPop(key string) string {
	return defStoreKeeper.SPop(key)
}

func SPopN(key string, count int64) []string {
	return defStoreKeeper.SPopN(key, count)
}

func SIsMember(key string, member interface{}) bool {
	return defStoreKeeper.SIsMember(key, member)
}

func ZAdd(key string, members ...*redis.Z) error {
	return defStoreKeeper.ZAdd(key, members...)
}
func ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	return defStoreKeeper.ZRangeWithScores(key, start, stop)
}
func ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	return defStoreKeeper.ZRevRangeWithScores(key, start, stop)
}

func ZRevRank(key, member string) (int64, error) {
	return defStoreKeeper.ZRevRank(key, member)
}

func ZRank(key, member string) (int64, error) {
	return defStoreKeeper.ZRank(key, member)
}
func ZScore(key, member string) float64 {
	return defStoreKeeper.ZScore(key, member)
}

func ZIncrBy(key string, increment float64, member string) (float64, error) {
	return defStoreKeeper.ZIncrBy(key, increment, member)
}

func ZRem(key string, members ...any) error {
	return defStoreKeeper.ZRem(key, members...)
}

func ZRemRangeByScore(key string, min, max string) error {
	return defStoreKeeper.ZRemRangeByScore(key, min, max)
}

func ZCard(key string) int64 {
	return defStoreKeeper.ZCard(key)
}

func Lock(key string) (*redsync.Mutex, context.Context, error) {
	return defStoreKeeper.Lock(key)
}

func Unlock(mutex *redsync.Mutex, ctx context.Context) error {
	return defStoreKeeper.Unlock(mutex, ctx)
}

func Incr(key string) (int64, error) {
	return defStoreKeeper.Incr(key)
}

func Decr(key string) (int64, error) {
	return defStoreKeeper.Decr(key)
}

func TxPipeline() redis.Pipeliner {
	return defStoreKeeper.TxPipeline()
}

func GetKey(key string) string {
	return defStoreKeeper.GetKey(key)
}
