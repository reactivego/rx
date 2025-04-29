package rx

func OnNext[T any](onNext func(T)) Pipe[T] {
	return Tap[T](func(next T, err error, done bool) {
		if !done {
			onNext(next)
		}
	})
}

func (observable Observable[T]) OnNext(f func(T)) Observable[T] {
	return OnNext(f)(observable)
}
