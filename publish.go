package rx

// Publish returns a multicasting Observable[T] for an underlying Observable[T] as a Connectable[T] type.
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
