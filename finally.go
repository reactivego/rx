package rx

func Finally[T any](f func(error)) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if done {
					f(err)
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Finally(f func(error)) Observable[T] {
	return Finally[T](f)(observable)
}
