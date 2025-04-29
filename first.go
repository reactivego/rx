package rx

func (observable Observable[T]) First(schedulers ...Scheduler) (value T, err error) {
	err = observable.Take(1).Assign(&value).Wait(schedulers...)
	return
}
