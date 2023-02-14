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
	"fmt"
	"testing"
	"time"
)

var hmacSampleSecret = "123456"

func TestJwtEncode(t *testing.T) {
	tokenString, exp, err := JwtEncode(1, 123456, 60*time.Second, hmacSampleSecret)
	fmt.Println(tokenString, exp, err)
}

func TestJwtDecode(t *testing.T) {
	// sample token string taken from the New example
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjowLCJ1c2VyX2lkIjozMDMxMjMsInN0YXJ0IjoxNjcwNTg1NjY1LCJleHBpcmUiOjE2NzA4NDQ4NjV9.SxtEq6oIKfV2rUReD0iVU150nryWg70c5DM3qnNPi3A"
	data, err := JwtDecode(tokenString, hmacSampleSecret)
	fmt.Println(data, err)
}
