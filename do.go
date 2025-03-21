package rx

func Do[T any](f func(T)) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done {
					f(next)
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Do(f func(next T)) Observable[T] {
	return Do[T](f)(observable)
}
