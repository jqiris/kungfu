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

package coroutines

import (
	"fmt"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	w := NewWaitCoroutines()
	w.AddCoroutine(func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("hello1")
	})
	w.AddCoroutine(func() {
		fmt.Println("hello2")
	})
	w.AddCoroutine(func() {
		fmt.Println("hello3")
	})
	w.AddCoroutine(func() {
		fmt.Println("hello4")
	})
	w.AddCoroutine(func() {
		fmt.Println("hello5")
	})
	w.Wait()
}
