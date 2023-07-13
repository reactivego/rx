package x

func CatchError[T any](selector func(err error, caught Observable[T]) Observable[T]) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done || err == nil {
					observe(next, err, done)
				} else {
					selector(err, observable)(observe, scheduler, subscriber)
				}
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) CatchError(selector func(err error, caught Observable[T]) Observable[T]) Observable[T] {
	return observable.Pipe(CatchError[T](selector))
}
