package x

func Send[T any](ch chan<- T) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			observable(func(next T, err error, done bool) {
				if !done {
					ch <- next
				}
				observe(next, err, done)
			}, scheduler, subscriber)
		}
	}
}
