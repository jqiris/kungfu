package coroutines

import (
	"sync"

	"github.com/jqiris/kungfu/v2/utils"
)

type CoroutineHandler func()

type WaitCoroutines struct {
	wg   *sync.WaitGroup
	list []CoroutineHandler
}

func NewWaitCoroutines() *WaitCoroutines {
	return &WaitCoroutines{
		wg:   new(sync.WaitGroup),
		list: make([]CoroutineHandler, 0),
	}
}

func (w *WaitCoroutines) AddCoroutine(handler CoroutineHandler) {
	w.list = append(w.list, handler)
}

func (w *WaitCoroutines) Wait() {
	n := len(w.list)
	w.wg.Add(n)
	for i := 0; i < len(w.list); i++ {
		handler := w.list[i]
		go utils.SafeRun(func() {
			defer w.wg.Done()
			handler()
		})
	}
	w.wg.Wait()
}
