package rx

import "context"

func (observable Observable[T]) Slice(ctx ...context.Context) (slice []T, err error) {
	err = observable.Append(&slice).Wait(ctx...)
	return
}
