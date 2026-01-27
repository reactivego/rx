package rx

import "context"

func (observable Observable[T]) Last(ctx ...context.Context) (value T, err error) {
	err = observable.Assign(&value).Wait(ctx...)
	return
}
