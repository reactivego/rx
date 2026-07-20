package rx

import "context"

// Go subscribes to the observable and starts execution on a separate goroutine.
// It ignores all emissions from the observable sequence, making it useful when
// you only care about side effects and not the actual values. By default, it
// uses the Goroutine scheduler, but a scheduler can be provided via context
// using SchedulerContextWith. Returns a Subscription that can be used to cancel
// the subscription when no longer needed.
func (observable Observable[T]) Go(ctx ...context.Context) Subscription {
	return observable.Subscribe(contextWithScheduler(ctx, Goroutine), Ignore[T]())
}
