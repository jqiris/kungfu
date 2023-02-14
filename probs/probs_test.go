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
