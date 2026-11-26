package windows

import "sync"

type Result struct {
	value any
	err   error
}

type ResultFuture struct {
	once sync.Once
	done chan Result
}

func newResultFuture() *ResultFuture {
	return &ResultFuture{
		done: make(chan Result, 1),
	}
}

func (f *ResultFuture) Complete(value any, err error) {
	f.once.Do(func() {
		f.done <- Result{
			value: value,
			err:   err,
		}
		close(f.done)
	})
}

func (f *ResultFuture) Await() (any, error) {
	result, ok := <-f.done
	if !ok {
		return nil, ErrWindowClosed
	}

	return result.value, result.err
}

func (f *ResultFuture) AwaitContext(ctxDone <-chan struct{}) (any, error) {
	select {
	case result, ok := <-f.done:
		if !ok {
			return nil, ErrWindowClosed
		}
		return result.value, result.err

	case <-ctxDone:
		return nil, contextCanceledError()
	}
}

func contextCanceledError() error {
	return errContextCanceled
}
