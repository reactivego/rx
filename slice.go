package x

func (observable Observable[T]) Slice(schedulers ...Scheduler) (slice []T) {
	observable.Subscribe(Collect(&slice), schedulers...).Wait()
	return
}
