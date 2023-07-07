package x

func (observable Observable[T]) Go(schedulers ...Scheduler) Subscription {
	if len(schedulers) == 0 {
		return observable.Subscribe(Ignore[T](), Goroutine)
	} else {
		return observable.Subscribe(Ignore[T](), schedulers[0])
	}
}
