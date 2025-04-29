package rx

func OnError[T any](onError func(error)) Pipe[T] {
	return Tap[T](func(next T, err error, done bool) {
		if done && err != nil {
			onError(err)
		}
	})
}

func (observable Observable[T]) OnError(f func(error)) Observable[T] {
	return OnError[T](f)(observable)
}
