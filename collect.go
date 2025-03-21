package rx

func Collect[T any](slice *[]T) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done {
					*slice = append(*slice, next)
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Collect(slice *[]T) Observable[T] {
	return Collect[T](slice)(observable)
}
