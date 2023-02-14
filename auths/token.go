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
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	TokenTypeNormal  = 0 //普通token
	TokenTypeRefresh = 1 //刷新token
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type JwtClaims struct {
	Type   int   `json:"type"`    //类型 0-普通tokn 1-刷新token
	UserId int64 `json:"user_id"` //用户ID
	Start  int64 `json:"start"`   //开始时间
	Expire int64 `json:"expire"`  //过期时间
}

func NewJwtClaims(typ int, userId int64, expire time.Duration) (*JwtClaims, time.Time) {
	start := time.Now()
	exp := start.Add(expire)
	return &JwtClaims{
		Type:   typ,
		UserId: userId,
		Start:  start.Unix(),
		Expire: exp.Unix(),
	}, exp
}

func (c *JwtClaims) Valid() error {
	if c.UserId < 1 {
		return ErrInvalidToken
	}
	nowUnix := time.Now().Unix()
	//还没开始
	if nowUnix < c.Start {
		return ErrInvalidToken
	}
	//已经过期
	if nowUnix > c.Expire {
		return ErrInvalidToken
	}
	return nil
}

func (c *JwtClaims) TypeRefresh() bool {
	return c.Type == TokenTypeRefresh
}

// jwt encode
func JwtEncode(typ int, userId int64, expire time.Duration, secret string) (string, time.Time, error) {
	claims, exp := NewJwtClaims(typ, userId, expire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, exp, err
}

// jwt decode
func JwtDecode(tokenString string, secret string) (*JwtClaims, error) {
	claims := &JwtClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	return claims, err
}
