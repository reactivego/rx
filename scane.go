package observable

func ScanE[T, U any](observable Observable[T], seed U, accumulator func(acc U, next T) (U, error)) Observable[U] {
	return func(observe Observer[U], scheduler Scheduler, subscriber Subscriber) {
		state := seed
		observable(func(next T, err error, done bool) {
			if !done {
				state, err = accumulator(state, next)
				observe(state, err, err != nil)
			} else {
				var zero U
				observe(zero, err, done)
			}
		}, scheduler, subscriber)
	}
}

func (observable Observable[T]) ScanE(seed T, accumulator func(acc, next T) (T, error)) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		state := seed
		observable(func(next T, err error, done bool) {
			if !done {
				state, err = accumulator(state, next)
				observe(state, err, err != nil)
			} else {
				var zero T
				observe(zero, err, done)
			}
		}, scheduler, subscriber)
	}
}
