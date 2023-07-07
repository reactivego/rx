package x

func (observable Observable[T]) Wait(schedulers ...Scheduler) error {
	if len(schedulers) == 0 {
		return observable.Subscribe(Ignore[T](), NewScheduler()).Wait()
	} else {
		return observable.Subscribe(Ignore[T](), schedulers[0]).Wait()
	}
}
