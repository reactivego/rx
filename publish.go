package rx

func (observable Observable[T]) Publish() Connectable[T] {
	observe, multicaster := Multicast[T](1)
	connector := func(scheduler Scheduler, subscriber Subscriber) {
		observer := func(next T, err error, done bool) {
			observe(next, err, done)
			if done {
				subscriber.Unsubscribe()
			}
		}
		observable(observer, scheduler, subscriber)
	}
	return Connectable[T]{Observable: multicaster, Connector: connector}
}
