package observable

func (observable Observable[T]) Take(n int) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		taken := 0
		observable(func(next T, err error, done bool) {
			if taken < n {
				observe(next, err, done)
				if !done {
					taken++
					if taken >= n {
						var zero T
						observe(zero, nil, true)
					}
				}
			}
		}, scheduler, subscriber)
	}
}
