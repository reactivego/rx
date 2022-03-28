package observable

func (observable Observable[T]) Subscribe(observe Observer[T], schedulers ...Scheduler) Subscription {
	if len(schedulers) == 0 {
		schedulers = []Scheduler{NewScheduler()}
	}
	scheduler := schedulers[0]
	subscription := newSubscription(scheduler.Wait)
	if scheduler.IsConcurrent() {
		subscription.wait = nil
	}
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
