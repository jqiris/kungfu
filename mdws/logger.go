/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package mdws

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jqiris/kungfu/v2/logger"
)

type LogFormatterParams struct {
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// BodySize is the size of the Response Body
	BodySize int
	//userId
	UserId int64
}
type ReqLogger struct {
	recorder    *logger.Logger
	identityKey string
}

func NewReqLogger(suffix string, identityKey string) *ReqLogger {
	return &ReqLogger{
		recorder:    logger.WithSuffix(suffix),
		identityKey: identityKey,
	}
}

func (r *ReqLogger) Record(item LogFormatterParams) {
	r.recorder.Infof("%d|%v|%s|%s|%s|%d", item.StatusCode, item.Latency, item.ClientIP, item.Method, item.Path, item.UserId)
}

func GinLogger(suffix string, identityKey string) gin.HandlerFunc {
	recorder := NewReqLogger(suffix, identityKey)
	return func(c *gin.Context) {
		//start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		//main exec
		c.Next()
		param := LogFormatterParams{}
		// Stop timer
		param.Latency = time.Since(start)
		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path
		param.UserId = c.GetInt64(recorder.identityKey)
		recorder.Record(param)
	}
}
