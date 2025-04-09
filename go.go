package rx

func (observable Observable[T]) Go() Subscription {
	return observable.Subscribe(Ignore[T](), Goroutine)
}
