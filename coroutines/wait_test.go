package coroutines

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	w := NewWaitCoroutines()
	w.AddCoroutine(func(wg *sync.WaitGroup) {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		fmt.Println("hello1")
	})
	w.AddCoroutine(func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("hello2")
	})
	w.AddCoroutine(func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("hello3")
	})
	w.AddCoroutine(func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("hello4")
	})
	w.AddCoroutine(func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("hello5")
	})
	w.Wait()
}
