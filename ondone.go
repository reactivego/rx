package rx

func OnDone[T any](onDone func(error)) Pipe[T] {
	return Tap[T](func(next T, err error, done bool) {
		if done {
			onDone(err)
		}
	})
}

func (observable Observable[T]) OnDone(f func(error)) Observable[T] {
	return OnDone[T](f)(observable)
}
