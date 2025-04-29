package rx

func (observable Observable[T]) Slice(schedulers ...Scheduler) (slice []T, err error) {
	err = observable.Append(&slice).Wait(schedulers...)
	return
}
