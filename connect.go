package x

// Connect instructs a ConnectableObservable[T] to subscribe to its source
// and begin emitting items to its subscribers. Connect accepts an optional
// scheduler argument.
func (connectable Connectable) Connect(schedulers ...Scheduler) Subscription {
	if len(schedulers) == 0 {
		schedulers = []Scheduler{NewScheduler()}
	}
	scheduler := schedulers[0]
	subscription := newSubscription(scheduler)
	connectable(schedulers[0], subscription)
	return subscription
}
