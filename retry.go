package rx

import "time"

func Retry[T any](limit ...int) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		backoff := func(int) time.Duration { return 1 * time.Millisecond }
		return observable.RetryTime(backoff, limit...)
	}
}

func (observable Observable[T]) Retry(limit ...int) Observable[T] {
	return Retry[T](limit...)(observable)
}
