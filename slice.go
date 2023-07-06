package x

func (observable Observable[T]) Slice(schedulers ...Scheduler) (slice []T, err error) {
	err = observable.Subscribe(Collect(&slice), schedulers...).Wait()
	return
}
