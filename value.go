package x

func (observable Observable[T]) Value(schedulers ...Scheduler) (value T, err error) {
	err = observable.Take(1).Assign(&value).Wait(schedulers...)
	return
}
