package rx

func (observable Observable[T]) Publish() Connectable[T] {
	observe, multicaster := Multicast[T](1)
	connect := func(scheduler Scheduler, subscriber Subscriber) {
		observable(observe, scheduler, subscriber)
	}
	return Connectable[T]{Observable: multicaster, Connect: connect}
}
