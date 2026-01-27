package rx

import "context"

func (observable Observable[T]) First(ctx ...context.Context) (value T, err error) {
	err = observable.Take(1).Assign(&value).Wait(ctx...)
	return
}
