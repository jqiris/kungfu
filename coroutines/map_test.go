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
	"math/rand"
	"testing"
	"time"

	"github.com/jqiris/kungfu/v2/jobs"
)

func TestNumberMap(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	data := NewNumberMap[int, int]()
	wait := make(chan struct{}, 2)
	jobs.AddJob(1*time.Second, jobs.NewJob(func() {
		for i := 0; i < 100; i++ {
			data.Store(i, rand.Intn(100))
		}
		fmt.Println("end1")
		wait <- struct{}{}
	}))
	jobs.AddJob(1*time.Second, jobs.NewJob(func() {
		for i := 100; i < 200; i++ {
			data.Incre(i, rand.Intn(200))
		}
		fmt.Println("end2")
		wait <- struct{}{}
	}))
	jobs.AddJob(1*time.Second, jobs.NewJob(func() {
		stop := make(chan struct{})
		num := 0
		for {
			select {
			case <-wait:
				num++
				if num == 2 {
					data.Range(func(k, v int) bool {
						data.Incre(k, rand.Intn(300))
						fmt.Println(k, v)
						return true
					})
					fmt.Println("range end")
					stop <- struct{}{}
				}
			case <-stop:
				return
			}
		}

	}))
	select {}
}
