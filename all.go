package rx

import "iter"

func (observable Observable[T]) All(scheduler ...Scheduler) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		index := -1
		yielder := func(value T) bool {
			index++
			return yield(index, value)
		}
		err := observable.TakeWhile(yielder).Wait(scheduler...)
		_ = err // ignores error! so will fail silently
	}
}
