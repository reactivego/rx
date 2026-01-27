package rx

import (
	"context"
	"iter"
)

func (observable Observable[T]) Values(ctx ...context.Context) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		// Wrap yield to match TakeWhile's expected signature
		yielding := func(value T) bool {
			return yield(value, nil)
		}
		err := observable.TakeWhile(yielding).Wait(ctx...)
		if err != nil {
			var zero T
			yield(zero, err)
		}
	}
}
