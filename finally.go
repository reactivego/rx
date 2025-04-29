package rx

// Finally is  shortcut for Tap(OnDone(f))
func Finally[T any](f func(error)) Pipe[T] {
	return Tap(OnDone[T](f))
}

// Finally is  shortcut for Tap(OnDone(f))
func (observable Observable[T]) Finally(f func(error)) Observable[T] {
	return Finally[T](f)(observable)
}
