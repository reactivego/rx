package rx

import "iter"

func (observable Observable[T]) Values(scheduler ...Scheduler) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		// Wrap yield to match TakeWhile's expected signature
		yielding := func(value T) bool {
			return yield(value, nil)
		}
		err := observable.TakeWhile(yielding).Wait(scheduler...)
		if err != nil {
			var zero T
			yield(zero, err)
		}
	}
}
