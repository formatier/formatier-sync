package future

import (
	"sync"
)

type FutureConfig[T any] struct {
	AsyncCallback bool
	AfterError    func(err error)
	AfterValue    func(v T)
}

func (fc *FutureConfig[T]) runAfterError(err error) {
	if fc.AfterError != nil {
		if fc.AsyncCallback {
			go fc.AfterError(err)
		} else {
			fc.AfterError(err)
		}
	}
}

func (fc *FutureConfig[T]) runAfterValue(v T) {
	if fc.AfterValue != nil {
		if fc.AsyncCallback {
			go fc.AfterValue(v)
		} else {
			fc.AfterValue(v)
		}
	}
}

func Warp[T any](callback func() (T, error), cfg *FutureConfig[T]) *Future[T] {
	if cfg == nil {
		cfg = &FutureConfig[T]{}
	}
	valueChan := make(chan T)
	errorChan := make(chan error)

	ft := &Future[T]{
		config:    cfg,
		valueChan: valueChan,
		errorChan: errorChan,
	}
	ft.initLazyValue()

	go func() {
		v, err := callback()
		if err != nil {
			ft.SendError(err)
			return
		}

		ft.SendValue(v)
	}()

	return ft
}

func NewFuture[T any](cfg *FutureConfig[T]) *Future[T] {
	if cfg == nil {
		cfg = &FutureConfig[T]{}
	}
	valueChan := make(chan T)
	errorChan := make(chan error)

	ft := &Future[T]{
		config:    cfg,
		valueChan: valueChan,
		errorChan: errorChan,
	}
	ft.initLazyValue()

	return ft
}

type Future[T any] struct {
	config    *FutureConfig[T]
	valueChan chan T
	errorChan chan error

	lazyValue func() (T, error)
}

func (f *Future[T]) SendValue(v T) {
	f.valueChan <- v
}

func (f *Future[T]) SendError(err error) {
	f.errorChan <- err
}

func (f *Future[T]) initLazyValue() {
	f.lazyValue = sync.OnceValues(func() (T, error) {
		defer close(f.errorChan)
		defer close(f.valueChan)

		select {
		case val := <-f.valueChan:
			f.config.runAfterValue(val)
			return val, nil
		case err := <-f.errorChan:
			f.config.runAfterError(err)
			var zero T
			return zero, err
		}
	})
}

func (f *Future[T]) Wait() (T, error) {
	return f.lazyValue()
}

func (f *Future[T]) WaitAsync(valueCallback func(v T), errorCallback func(err error)) {
	go func() {
		v, err := f.lazyValue()
		if err != nil {
			errorCallback(err)
			return
		}
		valueCallback(v)
	}()
}
