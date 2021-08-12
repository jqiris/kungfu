package tests

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var rule = regexp.MustCompile(`.*\/(.+\/.+\.go)`)

func TestEventType(t *testing.T) {
	a := "//gate/gate_2001"
	as := strings.Split(a, "/")

	fmt.Printf("%#v,len:%d \n", as, len(as))
	b := strings.Index(a, "/")
	fmt.Println(b)
}

func TestPathLocation(t *testing.T) {
	a := "H:/go/pkg/mod/github.com/jqiris/kungfu@v0.0.0-20210812091450-7f736d7f026f/rpcx/nats.go"
	b := pathA(a)
	fmt.Println(b)
	c := pathB(a)
	fmt.Println(c)
}
func pathA(a string) string {
	b := strings.Split(a, "/")
	lb := len(b)
	return b[lb-2] + "/" + b[lb-1]
}

func pathB(a string) string {
	res := rule.FindAllStringSubmatch(a, -1)
	if len(res) > 0 && len(res[0]) > 1 {
		return res[0][1]
	}
	return ""
}

func BenchmarkPatha(b *testing.B) {
	a := "H:/go/pkg/mod/github.com/jqiris/kungfu@v0.0.0-20210812091450-7f736d7f026f/rpcx/nats.go"
	for i := 0; i < b.N; i++ {
		pathA(a)
	}
}

func BenchmarkPathb(b *testing.B) {
	a := "H:/go/pkg/mod/github.com/jqiris/kungfu@v0.0.0-20210812091450-7f736d7f026f/rpcx/nats.go"
	for i := 0; i < b.N; i++ {
		pathB(a)
	}
}
