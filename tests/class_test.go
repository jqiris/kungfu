package tests

import (
	"fmt"
	"testing"
)

type Human struct {
	name  string
	age   int
	phone string
	call  func()
}

type Student struct {
	Human  //匿名字段
	school string
}

type Employee struct {
	Human   //匿名字段
	company string
}

//Human定义method
func (h *Human) SayHi() {
	fmt.Printf("Hi, I am %s you can call me on %s\n", h.name, h.phone)
}
func (h *Human) Call() {
	h.call()
}

func (h *Human) SetCall(handle func()) {
	h.call = handle
}

//Employee的method重写Human的method
func (e *Employee) SayHi() {
	fmt.Printf("Hi, I am %s, I work at %s. Call me on %s\n", e.name,
		e.company, e.phone) //Yes you can split into 2 lines here.
}

func TestCover(t *testing.T) {
	human := Human{"Mark", 25, "222-222-YYYY", nil}
	human.SetCall(func() { fmt.Println("Human call") })
	study := Employee{human, "mt"}
	study.SetCall(func() { fmt.Println("Employee call") })
	human.Call()
	study.Call()

}
