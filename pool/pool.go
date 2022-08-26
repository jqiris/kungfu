package pool

import (
	"errors"
	"sync"

	"github.com/jqiris/kungfu/v2/logger"
)

type Poolable interface {
	Close() error
}

type Pool[T Poolable] struct {
	lock      sync.Mutex
	resources chan T
	factory   func() (T, error)
	closed    bool
}

var ErrPoolClosed = errors.New("Pool has been closed")

// New 函数工厂，指定有缓冲通道大小
func New[T Poolable](fn func() (T, error), size uint) (*Pool[T], error) {
	if size <= 0 {
		return nil, errors.New("size value negative")
	}
	return &Pool[T]{
		factory:   fn,
		resources: make(chan T, size),
	}, nil
}

// Acquire 池中获取资源
func (p *Pool[T]) Acquire() (T, error) {
	var def T
	select {
	case r, ok := <-p.resources:
		if !ok {
			return def, ErrPoolClosed
		}
		return r, nil
	default:
		logger.Warn("Acquire: New Resource")
		return p.factory()
	}
}

// Release 池中释放资源
func (p *Pool[T]) Release(r T) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.closed {
		r.Close()
		return
	}

	select {
	case p.resources <- r: //放入队列
	default: //队列已满,则关闭
		logger.Warn("Release: Closing")
		r.Close()
	}
}

// Close 池关闭所有现有资源
func (p *Pool[T]) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.closed {
		return
	}

	p.closed = true
	close(p.resources) //清空通道资源前将通道关闭，否则会产生死锁
	for r := range p.resources {
		r.Close()
	}
}
