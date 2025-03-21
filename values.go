package x

import "iter"

func (observable Observable[T]) Values(scheduler ...Scheduler) iter.Seq[T] {
	return func(yield func(T) bool) {
		// ignores error! so will fail silently
		observable.TakeWhile(yield).Wait(scheduler...)
	}
}
