package rx

import "context"

// Connect instructs a Connectable[T] to subscribe to its source and begin
// emitting items to its subscribers. The scheduler can be provided via context
// using WithScheduler, otherwise a new serial scheduler is created.
func (connect Connector) Connect(ctx ...context.Context) Subscription {
	subscription := newSubscription(ctx...)
	connect(subscription.scheduler, subscription)
	return subscription
}
