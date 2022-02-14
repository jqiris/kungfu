package tests

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"testing"
	"time"
)

var rule = regexp.MustCompile(`.*\/(.+\/.+\.go)`)
var pathMaps = make(map[string]string)

func TestEventType(t *testing.T) {
	a := "//gate/gate_2001"
	as := strings.Split(a, "/")

	fmt.Printf("%#v,len:%d \n", as, len(as))
	b := strings.Index(a, "/")
	fmt.Println(b)
}

func TestPathLocation(t *testing.T) {
	a := "H:/go/pkg/mod/github.com/jqiris/kungfu@v0.0.0-20210812091450-7f736d7f026f/rpc/nats.go"
	b := pathA(a)
	fmt.Println(b)
	c := pathB(a)
	fmt.Println(c)
	d := pathC(a)
	fmt.Println(d)
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
func pathC(a string) string {
	if v, ok := pathMaps[a]; ok {
		return v
	}
	b := strings.Split(a, "/")
	lb := len(b)
	r := b[lb-2] + "/" + b[lb-1]
	pathMaps[a] = r
	return r
}

func BenchmarkPatha(b *testing.B) {
	a := "H:/go/pkg/mod/github.com/jqiris/kungfu@v0.0.0-20210812091450-7f736d7f026f/rpc/nats.go"
	for i := 0; i < b.N; i++ {
		pathA(a)
	}
}

func BenchmarkPathb(b *testing.B) {
	a := "H:/go/pkg/mod/github.com/jqiris/kungfu@v0.0.0-20210812091450-7f736d7f026f/rpc/nats.go"
	for i := 0; i < b.N; i++ {
		pathB(a)
	}
}
func BenchmarkPathc(b *testing.B) {
	a := "H:/go/pkg/mod/github.com/jqiris/kungfu@v0.0.0-20210812091450-7f736d7f026f/rpc/nats.go"
	for i := 0; i < b.N; i++ {
		pathC(a)
	}
}

func TestRand(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	a := int(rand.Float64() * 500)
	fmt.Println(a)
}
func abc(a any) {
	switch v := a.(type) {
	case int:
		fmt.Println("int:", v)
	case string:
		fmt.Println("string:", v)
	}
}
func TestGenericType(t *testing.T) {
	abc(31)
	abc("welcome")
}
