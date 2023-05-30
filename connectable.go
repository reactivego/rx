package x

// Connectable represents the subscription behavior of an observable that
// has different behavior for the source and the observable itself. Used
// with multicast observables to allow subscriptions to the observable
// before connecting (i.e. subscribing) to the source.
type Connectable func(Scheduler, Subscriber)

// Connect returns a subscription object that can be used to cancel the
// subscription and to wait for completion. It is used with connectables
// that do not begin emitting items until Connect is called. This allows
// for subscriptions to the observable to be established before connecting
// to the source which then starts emitting items. The subscription object
// returned by Connect can be used to cancel the subscription and to wait
// for completion of the observable. This is useful in scenarios where the
// observable is expensive to create or where the subscription needs to be
// established before the observable starts emitting items.
func (connectable Connectable) Connect(schedulers ...Scheduler) Subscription {
	if len(schedulers) == 0 {
		schedulers = []Scheduler{NewScheduler()}
	}
	scheduler := schedulers[0]
	subscription := newSubscription(scheduler)
	connectable(schedulers[0], subscription)
	return subscription
}
