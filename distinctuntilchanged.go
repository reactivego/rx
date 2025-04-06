package rx

func DistinctUntilChanged[T any](equal func(T, T) bool) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], subscribeOn Scheduler, subscriber Subscriber) {
			var seen struct {
				initialized bool
				value       T
			}
			observer := func(next T, err error, done bool) {
				if !done {
					if seen.initialized && equal(seen.value, next) {
						return // skip equal
					} else {
						seen.initialized = true
						seen.value = next
					}
				}
				observe(next, err, done)
			}
			observable(observer, subscribeOn, subscriber)
		}
	}
}

func (observable Observable[T]) DistinctUntilChanged(equal func(T, T) bool) Observable[T] {
	return DistinctUntilChanged(equal)(observable)
}
