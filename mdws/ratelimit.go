package mdws

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/spf13/viper"
	"go.uber.org/ratelimit"
)

type RateLimiter struct {
	rateLimiterLock *sync.RWMutex
	rateLimiterMap  map[int]ratelimit.Limiter
	switchName      string
	rateName        string
}

func NewRateLimiter(switchName, rateName string) *RateLimiter {
	return &RateLimiter{
		rateLimiterLock: new(sync.RWMutex),
		rateLimiterMap:  make(map[int]ratelimit.Limiter),
		switchName:      switchName,
		rateName:        rateName,
	}
}

func (r *RateLimiter) getLimiter(rate int) (ratelimit.Limiter, bool) {
	r.rateLimiterLock.RLock()
	defer r.rateLimiterLock.RUnlock()
	v, ok := r.rateLimiterMap[rate]
	return v, ok
}
func (r *RateLimiter) setLimiter(rate int) ratelimit.Limiter {
	r.rateLimiterLock.Lock()
	defer r.rateLimiterLock.Unlock()
	v := ratelimit.New(rate)
	r.rateLimiterMap[rate] = v
	return v
}

func (r *RateLimiter) Take() (bool, int, time.Time) {
	var takeTime time.Time
	isOpen, rate := viper.GetBool(r.switchName), viper.GetInt(r.rateName)
	if !isOpen {
		return false, rate, takeTime
	}
	v, ok := r.getLimiter(rate)
	if !ok {
		v = r.setLimiter(rate)
	}
	takeTime = v.Take()
	return true, rate, takeTime
}

func GinRateLimit(switchName, rateName string) gin.HandlerFunc {
	isOpen, rate := false, 0
	var before, after time.Time
	limiter := NewRateLimiter(switchName, rateName)
	return func(c *gin.Context) {
		before = time.Now()
		isOpen, rate, after = limiter.Take()
		if !isOpen {
			c.Next()
			return
		}
		logger.Warnf("rate limit rate: %v,cosume:%v", rate, after.Sub(before))
		c.Next()
	}
}
