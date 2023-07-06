package x

func (observable Observable[T]) Wait(schedulers ...Scheduler) error {
	return observable.Subscribe(Ignore[T](), schedulers...).Wait()
}
