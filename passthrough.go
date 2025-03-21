package rx

func Passthrough[T any]() Pipe[T] {
	// Pipe scope
	return func(observable Observable[T]) Observable[T] {
		// Operator scope
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			// Subscribe scope
			observable(func(next T, err error, done bool) {
				// Observe scope
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Passthrough() Observable[T] {
	return Passthrough[T]()(observable)
}
