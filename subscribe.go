package rx

import "context"

// Subscribe subscribes to the observable with the given context, observer, and scheduler.
// The context will be available throughout the observable chain via subscriber.Context().
// Source operators (like Create, Recv) are responsible for checking context cancellation.
func (observable Observable[T]) Subscribe(ctx context.Context, observe Observer[T]) Subscription {
	subscription := newSubscription(ctx)
	observer := func(next T, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			var zero T
			observe(zero, err, true)
			subscription.done(err)
		}
	}
	observable(observer, subscription.scheduler, subscription)
	return subscription
}
