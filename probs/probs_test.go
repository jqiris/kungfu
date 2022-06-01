package probs

import (
	"fmt"
	"testing"
)

func TestProbs(t *testing.T) {
	data := map[string]int{
		"a": 10,
		"b": 20,
		"c": 70,
	}
	res := make(map[string]int)
	wgRand := NewWgRand(false)
	for k, v := range data {
		wgRand.AddElement(k, v)
	}
	for i := 0; i < 100000; i++ {
		choice, err := wgRand.GetRandomChoice()
		if err != nil {
			fmt.Println(err)
		} else {
			res[choice.(string)]++
		}
	}
	fmt.Println(res)
}
