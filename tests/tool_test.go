package tests

import (
	"fmt"
	"strings"
	"testing"
)

func TestEventType(t *testing.T) {
	a := "//gate/gate_2001"
	as := strings.Split(a, "/")

	fmt.Printf("%#v,len:%d \n", as, len(as))
	b := strings.Index(a, "/")
	fmt.Println(b)
}
