package rx

// Do is a shortcut for Tap(OnNext(f))
func Do[T any](f func(T)) Pipe[T] {
	return Tap(OnNext(f))
}

// Do is a shortcut for Tap(OnNext(f))
func (observable Observable[T]) Do(f func(next T)) Observable[T] {
	return Do(f)(observable)
}
