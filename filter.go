package rx

func Filter[T any](predicate func(T) bool) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if done || predicate(next) {
					observe(next, err, done)
				}
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Filter(predicate func(T) bool) Observable[T] {
	return Filter[T](predicate)(observable)
}
