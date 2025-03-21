package rx

func (observable Observable[T]) SubscribeOn(scheduler ConcurrentScheduler) Observable[T] {
	return func(observe Observer[T], _ Scheduler, subscriber Subscriber) {
		observable(observe, scheduler, subscriber)
	}
}
