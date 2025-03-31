package rx

import "iter"

func (observable Observable[T]) Values(scheduler ...Scheduler) iter.Seq[T] {
	return func(yield func(T) bool) {
		err := observable.TakeWhile(yield).Wait(scheduler...)
		_ = err // ignores error! so will fail silently
	}
}
