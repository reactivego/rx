package rx

func (observable Observable[T]) Subscribe(observe Observer[T], scheduler Scheduler) Subscription {
	subscription := newSubscription(scheduler)
	observer := func(next T, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			var zero T
			observe(zero, err, true)
			subscription.done(err)
		}
	}
	observable(observer, scheduler, subscription)
	return subscription
}
