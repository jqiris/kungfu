package coroutines

import "sync"

type CoroutineHandler func(wg *sync.WaitGroup)

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
	for _, handler := range w.list {
		go handler(w.wg)
	}
	w.wg.Wait()
}
