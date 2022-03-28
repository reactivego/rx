package observable

func Scan[T, U any](observable Observable[T], seed U, accumulator func(acc U, next T) U) Observable[U] {
	return func(observe Observer[U], scheduler Scheduler, subscriber Subscriber) {
		state := seed
		observable(func(next T, err error, done bool) {
			if !done {
				state = accumulator(state, next)
				observe(state, nil, false)
			} else {
				var zero U
				observe(zero, err, done)
			}
		}, scheduler, subscriber)
	}
}

func (observable Observable[T]) Scan(seed T, accumulator func(acc, next T) T) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		state := seed
		observable(func(next T, err error, done bool) {
			if !done {
				state = accumulator(state, next)
				observe(state, nil, false)
			} else {
				var zero T
				observe(zero, err, done)
			}
		}, scheduler, subscriber)
	}
}
