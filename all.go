package rx

import "iter"

// All2 converts an Observable of Tuple2 pairs into an iterator sequence. It
// transforms an Observable[Tuple2[T, U]] into an iter.Seq2[T, U] that yields
// each tuple's components. The scheduler parameter is optional and determines
// the execution context. Note: This method ignores any errors from the
// observable stream.
func All2[T, U any](observable Observable[Tuple2[T, U]], scheduler ...Scheduler) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		yielding := func(next Tuple2[T, U]) bool {
			return yield(next.First, next.Second)
		}
		err := observable.TakeWhile(yielding).Wait(scheduler...)
		_ = err // ignores error! so will fail silently
	}
}

// All converts an Observable stream into an iterator sequence that pairs each
// element with its index. It returns an iter.Seq2[int, T] which yields each
// element along with its position in the sequence. The scheduler parameter is
// optional and determines the execution context. Note: This method ignores any
// errors from the observable stream.
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
