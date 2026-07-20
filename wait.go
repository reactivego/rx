package rx

import (
	"context"
)

func (observable Observable[T]) Wait(ctx ...context.Context) error {
	return observable.Subscribe(contextWithScheduler(ctx), Ignore[T]()).Wait()
}
