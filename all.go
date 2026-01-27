package rx

import (
	"context"
	"iter"
)

// All converts an Observable stream into an iterator sequence that pairs each
// element with its index. It returns an iter.Seq2[Tuple2[int, T], error] which
// yields each element along with its position in the sequence as a Tuple2.
func (observable Observable[T]) All(ctx ...context.Context) iter.Seq2[Tuple2[int, T], error] {
	return func(yield func(Tuple2[int, T], error) bool) {
		index := -1
		yielding := func(next T) bool {
			index++
			return yield(Tuple2[int, T]{First: index, Second: next}, nil)
		}
		err := observable.TakeWhile(yielding).Wait(ctx...)
		if err != nil {
			var zero Tuple2[int, T]
			yield(zero, err)
		}
	}
}
