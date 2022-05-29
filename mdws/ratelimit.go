package mdws

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/spf13/viper"
	"go.uber.org/ratelimit"
)

var (
	rateLimiterMap = make(map[int]ratelimit.Limiter)
)

func GinRateLimit(switchName, rateName string) gin.HandlerFunc {
	var before, after time.Time
	isOpen, rate := false, 0
	return func(c *gin.Context) {
		isOpen, rate = viper.GetBool(switchName), viper.GetInt(rateName)
		if !isOpen {
			c.Next()
			return
		}
		v, ok := rateLimiterMap[rate]
		if !ok {
			v = ratelimit.New(rate)
			rateLimiterMap[rate] = v
		}
		before = time.Now()
		after = v.Take()
		logger.Infof("rate limit rate: %v,cosume:%v", rate, after.Sub(before))
		c.Next()
	}
}
