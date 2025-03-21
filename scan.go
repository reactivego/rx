package rx

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
