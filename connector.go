package rx

// Connector provides the Connect method for a Connectable[T].
type Connector func(Scheduler, Subscriber)

// Connect instructs a Connectable[T] to subscribe to its source and begin
// emitting items to its subscribers. Connect accepts an optional scheduler
// argument.
func (connect Connector) Connect(schedulers ...Scheduler) Subscription {
	if len(schedulers) == 0 {
		schedulers = []Scheduler{NewScheduler()}
	}
	scheduler := schedulers[0]
	subscription := newSubscription(scheduler)
	connect(schedulers[0], subscription)
	return subscription
}
