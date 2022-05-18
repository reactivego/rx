package x

func (observable Observable[T]) Value(schedulers ...Scheduler) (value T) {
	observable.Take(1).Subscribe(Assign(&value), schedulers...).Wait()
	return
}
