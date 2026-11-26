package windows

import (
	"context"
	"sync"
)

type Result struct {
	value any
	err   error
}

// ResultFuture is a broadcast-style one-shot result. Unlike a buffered result
// channel, any number of Await calls can observe the same completed value.
type ResultFuture struct {
	once   sync.Once
	done   chan struct{}
	mu     sync.RWMutex
	result Result
}

func newResultFuture() *ResultFuture {
	return &ResultFuture{done: make(chan struct{})}
}

func (f *ResultFuture) Complete(value any, err error) {
	f.once.Do(func() {
		f.mu.Lock()
		f.result = Result{value: value, err: err}
		f.mu.Unlock()
		close(f.done)
	})
}

func (f *ResultFuture) Await() (any, error) {
	<-f.done
	return f.value()
}

func (f *ResultFuture) AwaitContext(ctx context.Context) (any, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-f.done:
		return f.value()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (f *ResultFuture) value() (any, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.result.value, f.result.err
}
