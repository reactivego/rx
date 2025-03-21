package x

func Take[T any](n int) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			taken := 0
			observable(func(next T, err error, done bool) {
				if done || taken < n {
					observe(next, err, done)
				} else {
					var zero T
					observe(zero, nil, true)
				}
				taken++
			}, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) Take(n int) Observable[T] {
	return Take[T](n)(observable)
}
