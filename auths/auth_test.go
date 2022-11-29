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
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoxLCJ1c2VyX2lkIjoxMjM0NTYsInN0YXJ0IjoxNjY5NjkyMzM3LCJleHBpcmUiOjE2Njk2OTIzOTd9.5RRcwdlNfrMnYFYiucu9OFEZ-GMeBJFVAD3PIjGvwqE"
	data, err := JwtDecode(tokenString, hmacSampleSecret)
	fmt.Println(data, err)
}
