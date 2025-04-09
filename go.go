package rx

// Go subscribes to the observable on a new goroutine, ignoring all emissions.
// This method provides a convenient way to start consuming the observable asynchronously
// without processing its values. Returns a Subscription which can be used to unsubscribe.
func (observable Observable[T]) Go() Subscription {
	return observable.Subscribe(Ignore[T](), Goroutine)
}
