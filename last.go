package rx

func (observable Observable[T]) Last(schedulers ...Scheduler) (value T, err error) {
	err = observable.Assign(&value).Wait(schedulers...)
	return
}
