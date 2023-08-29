package x

func Catch[T any](other Observable[T]) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done || err == nil {
					observe(next, err, done)
				} else {
					other(observe, scheduler, subscriber)
				}
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Catch(other Observable[T]) Observable[T] {
	return Catch[T](other)(observable)
}
