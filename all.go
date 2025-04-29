package rx

import "iter"

func All2[T, U any](observable Observable[Tuple2[T, U]], scheduler ...Scheduler) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		yielding := func(next Tuple2[T, U]) bool {
			return yield(next.First, next.Second)
		}
		err := observable.TakeWhile(yielding).Wait(scheduler...)
		_ = err // ignores error! so will fail silently
	}
}

func (observable Observable[T]) All(scheduler ...Scheduler) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		index := -1
		yielding := func(value T) bool {
			index++
			return yield(index, value)
		}
		err := observable.TakeWhile(yielding).Wait(scheduler...)
		_ = err // ignores error! so will fail silently
	}
}
