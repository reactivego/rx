package rx

func AutoUnsubscribe[T any]() Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			subscriber = subscriber.Add()
			observer := func(next T, err error, done bool) {
				observe(next, err, done)
				if done {
					subscriber.Unsubscribe()
				}
			}
			observable(observer, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) AutoUnsubscribe() Observable[T] {
	return AutoUnsubscribe[T]()(observable)
}
