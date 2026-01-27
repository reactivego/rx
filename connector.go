package rx

import "context"

// Connector provides the Connect method for a Connectable[T].
type Connector func(Scheduler, Subscriber)

// Connect instructs a Connectable[T] to subscribe to its source and begin
// emitting items to its subscribers. The scheduler can be provided via context
// using WithScheduler, otherwise a new serial scheduler is created.
func (connect Connector) Connect(ctx ...context.Context) Subscription {
	ctx = append(ctx, nil)
	subscription := newSubscription(ctx[0])
	connect(subscription.scheduler, subscription)
	return subscription
}
