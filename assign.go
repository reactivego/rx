package x

func Assign[T any](value *T) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done {
					*value = next
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}
