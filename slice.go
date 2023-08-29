package x

func (observable Observable[T]) Slice(schedulers ...Scheduler) (slice []T, err error) {
	err = observable.Collect(&slice).Wait(schedulers...)
	return
}
