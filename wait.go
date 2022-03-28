package observable

func (observable Observable[T]) Wait(schedulers ...Scheduler) error {
	return observable.Subscribe(Wait[T](), schedulers...).Wait()
}

func Wait[T any]() Observer[T] {
	return func(next T, err error, done bool) {}
}
