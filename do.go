package rx

func Do[T any](do func(T)) Pipe[T] {
	return Tap[T](func(next T, err error, done bool) {
		if !done {
			do(next)
		}
	})
}

func (observable Observable[T]) Do(f func(T)) Observable[T] {
	return Do(f)(observable)
}
