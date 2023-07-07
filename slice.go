package x

func (observable Observable[T]) Slice(schedulers ...Scheduler) (slice []T, err error) {
	err = observable.Pipe(Collect(&slice)).Wait(schedulers...)
	return
}
