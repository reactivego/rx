package x

import (
	"time"
)

func (observable Observable[T]) Retry(limit ...int) Observable[T] {
	backoff := func(int) time.Duration { return 1 * time.Millisecond }
	return observable.RetryTime(backoff, limit...)
}
