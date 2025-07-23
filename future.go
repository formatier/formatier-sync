package future

import (
	"context"
	"errors"
	"time"
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

func WarpErrorFuture[T any](callback func() (T, error), cfg *FutureConfig[T]) *ErrorFuture[T] {
	valueChan := make(chan T)
	errorChan := make(chan error)

	return &ErrorFuture[T]{
		config:    cfg,
		valueChan: valueChan,
		errorChan: errorChan,

		lazyValue: callback,
	}
}

func NewErrorFuture[T any](cfg *FutureConfig[T]) *ErrorFuture[T] {
	if cfg == nil {
		cfg = &FutureConfig[T]{}
	}
	valueChan := make(chan T)
	errorChan := make(chan error)

	return &ErrorFuture[T]{
		config:    cfg,
		valueChan: valueChan,
		errorChan: errorChan,
	}
}

func NewErrorFutureWithChan[T any](valueChan chan T, errorChan chan error, cfg *FutureConfig[T]) *ErrorFuture[T] {
	if cfg == nil {
		cfg = &FutureConfig[T]{}
	}
	return &ErrorFuture[T]{
		config:    cfg,
		valueChan: valueChan,
		errorChan: errorChan,
	}
}

type ErrorFuture[T any] struct {
	config    *FutureConfig[T]
	valueChan chan T
	errorChan chan error

	lazyValue func() (T, error)
}

func (f *ErrorFuture[T]) SendValue(v T) {
	f.valueChan <- v
}

func (f *ErrorFuture[T]) SendError(err error) {
	f.errorChan <- err
}

func (f *ErrorFuture[T]) runIfLazy() {
	if f.lazyValue != nil {
		v, err := f.lazyValue()
		if err != nil {
			f.SendError(err)
			return
		}
		f.SendValue(v)
	}
}

func (f *ErrorFuture[T]) Wait() (T, error) {
	defer close(f.errorChan)
	defer close(f.valueChan)

	go f.runIfLazy()

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		return val, nil
	case err := <-f.errorChan:
		f.config.runAfterError(err)
		var zero T
		return zero, err
	}
}

func (f *ErrorFuture[T]) WaitAsync(callback func(v T), errorCallback func(err error)) {
	defer close(f.errorChan)
	defer close(f.valueChan)

	go f.runIfLazy()

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		go callback(val)
	case err := <-f.errorChan:
		f.config.runAfterError(err)
		go errorCallback(err)
	}
}

func (f *ErrorFuture[T]) WaitCtx(ctx context.Context) (T, error) {
	defer close(f.errorChan)
	defer close(f.valueChan)

	go f.runIfLazy()

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		return val, nil
	case err := <-f.errorChan:
		f.config.runAfterError(err)
		var zero T
		return zero, err
	case <-ctx.Done():
		err := errors.New("context canceled")
		f.config.runAfterError(err)
		var zero T
		return zero, err
	}
}

func (f *ErrorFuture[T]) WaitCtxAsync(ctx context.Context, callback func(v T), errorCallback func(err error)) {
	defer close(f.errorChan)
	defer close(f.valueChan)

	go f.runIfLazy()

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		go callback(val)
	case err := <-f.errorChan:
		f.config.runAfterError(err)
		go errorCallback(err)
	case <-ctx.Done():
		err := errors.New("context canceled")
		f.config.runAfterError(err)
		go errorCallback(err)
	}
}

func (f *ErrorFuture[T]) WaitTimeout(duration time.Duration) (T, error) {
	defer close(f.errorChan)
	defer close(f.valueChan)

	go f.runIfLazy()

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		return val, nil
	case err := <-f.errorChan:
		f.config.runAfterError(err)
		var zero T
		return zero, err
	case <-time.After(duration):
		err := errors.New("timeout")
		f.config.runAfterError(err)
		var zero T
		return zero, err
	}
}

func (f *ErrorFuture[T]) WaitTimeoutAsync(duration time.Duration, callback func(v T), errorCallback func(err error)) {
	defer close(f.errorChan)
	defer close(f.valueChan)

	go f.runIfLazy()

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		go callback(val)
	case err := <-f.errorChan:
		f.config.runAfterError(err)
		go errorCallback(err)
	case <-time.After(duration):
		err := errors.New("timeout")
		f.config.runAfterError(err)
		go errorCallback(err)
	}
}

func (f *ErrorFuture[T]) WaitTimeoutCtx(ctx context.Context, duration time.Duration) (T, error) {
	defer close(f.errorChan)
	defer close(f.valueChan)

	go f.runIfLazy()

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		return val, nil
	case err := <-f.errorChan:
		f.config.runAfterError(err)
		var zero T
		return zero, err
	case <-time.After(duration):
		err := errors.New("timeout")
		f.config.runAfterError(err)
		var zero T
		return zero, err
	case <-ctx.Done():
		err := errors.New("context canceled")
		f.config.runAfterError(err)
		var zero T
		return zero, err
	}
}

func (f *ErrorFuture[T]) WaitTimeoutCtxAsync(ctx context.Context, duration time.Duration, callback func(v T), errorCallback func(err error)) {
	defer close(f.errorChan)
	defer close(f.valueChan)

	go f.runIfLazy()

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		go callback(val)
	case err := <-f.errorChan:
		f.config.runAfterError(err)
		go errorCallback(err)
	case <-time.After(duration):
		err := errors.New("timeout")
		f.config.runAfterError(err)
		go errorCallback(err)
	case <-ctx.Done():
		err := errors.New("context canceled")
		f.config.runAfterError(err)
		go errorCallback(err)
	}
}

func NewFuture[T any](cfg *FutureConfig[T]) *Future[T] {
	if cfg == nil {
		cfg = &FutureConfig[T]{}
	}
	valueChan := make(chan T)

	return &Future[T]{
		config:    cfg,
		valueChan: valueChan,
	}
}

func NewFutureWithChan[T any](valueChan chan T, cfg *FutureConfig[T]) *ErrorFuture[T] {
	if cfg == nil {
		cfg = &FutureConfig[T]{}
	}
	return &ErrorFuture[T]{
		config:    cfg,
		valueChan: valueChan,
	}
}

type Future[T any] struct {
	config    *FutureConfig[T]
	valueChan chan T
}

func (f *Future[T]) SendValue(v T) {
	f.valueChan <- v
}

func (f *Future[T]) Wait() T {
	defer close(f.valueChan)

	val := <-f.valueChan
	f.config.runAfterValue(val)
	return val
}

func (f *Future[T]) WaitAsync(callback func(v T)) {
	defer close(f.valueChan)

	val := <-f.valueChan
	f.config.runAfterValue(val)
	go callback(val)
}

func (f *Future[T]) WaitCtx(ctx context.Context) (T, error) {
	defer close(f.valueChan)

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		return val, nil
	case <-ctx.Done():
		err := errors.New("context canceled")
		f.config.runAfterError(err)
		var zero T
		return zero, err
	}
}

func (f *Future[T]) WaitCtxAsync(ctx context.Context, callback func(v T), errorCallback func(err error)) {
	defer close(f.valueChan)

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		go callback(val)
	case <-ctx.Done():
		err := errors.New("context canceled")
		f.config.runAfterError(err)
		go errorCallback(err)
	}
}

func (f *Future[T]) WaitTimeout(duration time.Duration) (T, error) {
	defer close(f.valueChan)

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		return val, nil
	case <-time.After(duration):
		err := errors.New("timeout")
		f.config.runAfterError(err)
		var zero T
		return zero, err
	}
}

func (f *Future[T]) WaitTimeoutAsync(duration time.Duration, callback func(v T), errorCallback func(err error)) {
	defer close(f.valueChan)

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		go callback(val)
	case <-time.After(duration):
		err := errors.New("timeout")
		f.config.runAfterError(err)
		go errorCallback(err)
	}
}

func (f *Future[T]) WaitTimeoutCtx(ctx context.Context, duration time.Duration) (T, error) {
	defer close(f.valueChan)

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		return val, nil
	case <-time.After(duration):
		err := errors.New("timeout")
		f.config.runAfterError(err)
		var zero T
		return zero, err
	case <-ctx.Done():
		err := errors.New("context canceled")
		f.config.runAfterError(err)
		var zero T
		return zero, err
	}
}

func (f *Future[T]) WaitTimeoutCtxAsync(ctx context.Context, duration time.Duration, callback func(v T), errorCallback func(err error)) {
	defer close(f.valueChan)

	select {
	case val := <-f.valueChan:
		f.config.runAfterValue(val)
		go callback(val)
	case <-time.After(duration):
		err := errors.New("timeout")
		f.config.runAfterError(err)
		go errorCallback(err)
	case <-ctx.Done():
		err := errors.New("context canceled")
		f.config.runAfterError(err)
		go errorCallback(err)
	}
}
