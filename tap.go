package rx

func Tap[T any](tap Observer[T]) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				tap(next, err, done)
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Tap(tap Observer[T]) Observable[T] {
	return Tap[T](tap)(observable)
}
