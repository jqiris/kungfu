package coroutines

import (
	"sync"

	"github.com/jqiris/kungfu/v2/utils"
)

type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64 | ~uint | ~uint32 | ~uint64 | ~complex64 | ~complex128
}

type NumberMap[K comparable, V Number] struct {
	lock *sync.RWMutex
	data map[K]V
}

func NewNumberMap[k comparable, v Number]() *NumberMap[k, v] {
	return &NumberMap[k, v]{
		lock: new(sync.RWMutex),
		data: make(map[k]v),
	}
}

func (m *NumberMap[K, V]) MarshalJSON() ([]byte, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return utils.JsonMarshal(m.data)
}

func (m *NumberMap[K, V]) UnmarshalJSON(data []byte) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return utils.JsonUnmarshal(data, &m.data)
}

func (m *NumberMap[K, V]) getLock(args ...bool) bool {
	lock := true
	if len(args) > 0 {
		lock = args[0]
	}
	return lock
}

func (m *NumberMap[K, V]) Load(k K, args ...bool) V {
	if m.getLock(args...) {
		m.lock.RLock()
		defer m.lock.RUnlock()
	}
	return m.data[k]
}

func (m *NumberMap[K, V]) LoadOk(k K, args ...bool) (V, bool) {
	if m.getLock(args...) {
		m.lock.RLock()
		defer m.lock.RUnlock()
	}
	v, ok := m.data[k]
	return v, ok
}

func (m *NumberMap[K, V]) Store(k K, v V, args ...bool) {
	if m.getLock(args...) {
		m.lock.Lock()
		defer m.lock.Unlock()
	}
	m.data[k] = v
}

func (m *NumberMap[K, V]) Incre(k K, increment V, args ...bool) {
	if m.getLock(args...) {
		m.lock.Lock()
		defer m.lock.Unlock()
	}
	m.data[k] += increment
}
func (m *NumberMap[K, V]) IncreOne(k K, args ...bool) {
	if m.getLock(args...) {
		m.lock.Lock()
		defer m.lock.Unlock()
	}
	m.data[k] += 1
}

func (m *NumberMap[K, V]) Decre(k K, decrement V, args ...bool) {
	if m.getLock(args...) {
		m.lock.Lock()
		defer m.lock.Unlock()
	}
	m.data[k] -= decrement
}

func (m *NumberMap[K, V]) DecreOne(k K, args ...bool) {
	if m.getLock(args...) {
		m.lock.Lock()
		defer m.lock.Unlock()
	}
	m.data[k] -= 1
}

func (m *NumberMap[K, V]) Range(visit func(k K, v V) bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for k, v := range m.data {
		if !visit(k, v) {
			break
		}
	}
}
