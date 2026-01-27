package rx

import (
	"context"
)

func (observable Observable[T]) Wait(ctx ...context.Context) error {
	return observable.Subscribe(withScheduler(ctx), Ignore[T]()).Wait()
}
