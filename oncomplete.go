package rx

func OnComplete[T any](onComplete func()) Pipe[T] {
	return Tap[T](func(next T, err error, done bool) {
		if done && err == nil {
			onComplete()
		}
	})
}

func (observable Observable[T]) OnComplete(f func()) Observable[T] {
	return OnComplete[T](f)(observable)
}
