package x

func (observable Observable[T]) Tap(tap func(T)) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		observable(func(next T, err error, done bool) {
			if !done {
				tap(next)
			}
			observe(next, err, done)
		}, scheduler, subscriber)
	}
}

func Tap[T any](tap func(T)) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done {
					tap(next)
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}
