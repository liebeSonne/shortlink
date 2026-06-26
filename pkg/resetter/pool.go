package resetter

import (
	"errors"
	"sync"
)

var InvalidPoolConstructorFuncErr = errors.New("invalid pool constructor function")

type Resetter interface {
	Reset()
}

type Pool[T Resetter] struct {
	pool sync.Pool
}

func (p *Pool[T]) Put(obj T) {
	obj.Reset()
	p.pool.Put(obj)
}
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

func NewPool[T Resetter](newFunc func() T) (*Pool[T], error) {
	if newFunc == nil {
		return nil, InvalidPoolConstructorFuncErr
	}
	return &Pool[T]{
		pool: sync.Pool{
			New: func() interface{} { return newFunc() },
		},
	}, nil
}
