package observable

func (observable Observable[T]) AutoUnsubscribe() Observable[T] {
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
