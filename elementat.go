package rx

func ElementAt[T any](n int) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			i := 0
			observer := func(next T, err error, done bool) {
				if !done {
					if i >= n {
						if i == n {
							observe(next, nil, false)
						} else {
							var zero T
							observe(zero, nil, true)
						}
					}
					i++
				} else {
					observe(next, err, done)
				}
			}
			observable(observer, scheduler, subscriber)
		}
	}
}

func (observable Observable[T]) ElementAt(n int) Observable[T] {
	return ElementAt[T](n)(observable)
}
