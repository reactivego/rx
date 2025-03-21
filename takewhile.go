package rx

func TakeWhile[T any](condition func(T) bool) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if done || condition(next) {
					observe(next, err, done)
				} else {
					var zero T
					observe(zero, nil, true)
				}
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) TakeWhile(condition func(T) bool) Observable[T] {
	return TakeWhile[T](condition)(observable)
}
