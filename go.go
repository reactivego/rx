package rx

// Go subscribes to the observable and starts execution on a separate goroutine.
// It ignores all emissions from the observable sequence, making it useful when you
// only care about side effects and not the actual values. By default, it uses the
// Goroutine scheduler, but an optional scheduler can be provided. Returns a
// Subscription that can be used to cancel the subscription when no longer needed.
func (observable Observable[T]) Go(schedulers ...Scheduler) Subscription {
	if len(schedulers) == 0 {
		return observable.Subscribe(Ignore[T](), Goroutine)
	} else {
		return observable.Subscribe(Ignore[T](), schedulers[0])
	}
}
