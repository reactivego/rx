package rx

func (observable Observable[T]) Publish() Connectable[T] {
	observe, multicaster := Multicast[T](1)
	connector := func(scheduler Scheduler, subscriber Subscriber) {
		observable(observe, scheduler, subscriber)
	}
	return Connectable[T]{Observable: multicaster, Connector: connector}
}
