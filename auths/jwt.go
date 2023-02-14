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

package auths

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/utils"
	"github.com/spf13/viper"
)

type JwtData interface {
	UserId() int64
}

type JwtUser interface {
	Authenticator(c *gin.Context) (interface{}, error)
	// User can define own Unauthorized func.
	Unauthorized(c *gin.Context, code int, message string)
	// User can define own LoginResponse func.
	LoginResponse(c *gin.Context, code int, message string, time time.Time)
	// User can define own LogoutResponse func.
	LogoutResponse(c *gin.Context, code int)
	// User can define own RefreshResponse func.
	RefreshResponse(c *gin.Context, code int, message string, time time.Time)
}

type JwtMiddleware struct {
	user JwtUser
}

func NewJwtMiddleware(user JwtUser) *JwtMiddleware {
	return &JwtMiddleware{
		user: user,
	}
}

func (g *JwtMiddleware) GetAuthor() *jwt.GinJWTMiddleware {
	cfg := viper.GetStringMapString("jwt")
	timeoutMinutes, maxRefreshMinutes := utils.StringToInt64(cfg["timeout_minutes"]), utils.StringToInt64(cfg["max_refresh_minutes"])
	author, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           cfg["realm"],
		Key:             []byte(cfg["key"]),
		Timeout:         time.Duration(timeoutMinutes) * time.Minute,
		MaxRefresh:      time.Duration(maxRefreshMinutes) * time.Minute,
		Authenticator:   g.user.Authenticator,
		Authorizator:    g.Verifier,
		PayloadFunc:     g.PayloadFunc,
		Unauthorized:    g.user.Unauthorized,
		LoginResponse:   g.user.LoginResponse,
		LogoutResponse:  g.user.LogoutResponse,
		RefreshResponse: g.user.RefreshResponse,
		IdentityHandler: g.IdentityHandler,
		IdentityKey:     cfg["identity_key"],
		TokenLookup:     "header: Authorization, query: token",
		TokenHeadName:   cfg["token_head_name"],
		TimeFunc:        time.Now,
	})
	if err != nil {
		logger.Fatal(err)
	}
	if err = author.MiddlewareInit(); err != nil {
		logger.Fatal(err)
	}
	return author
}

func (g *JwtMiddleware) PayloadFunc(data interface{}) jwt.MapClaims {
	identityKey := viper.GetString("jwt.identity_key")
	if v, ok := data.(JwtData); ok {
		return jwt.MapClaims{
			identityKey: v.UserId(),
		}
	}
	return jwt.MapClaims{}
}

func (g *JwtMiddleware) IdentityHandler(c *gin.Context) interface{} {
	identityKey := viper.GetString("jwt.identity_key")
	claims := jwt.ExtractClaims(c)
	return int64(claims[identityKey].(float64))
}

func (g *JwtMiddleware) Verifier(data interface{}, c *gin.Context) bool {
	if uid, ok := data.(int64); ok && uid > 0 {
		return true
	}
	return false
}
