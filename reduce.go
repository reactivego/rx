package observable

func Reduce[T, U any](observable Observable[T], seed U, accumulator func(acc U, next T) U) Observable[U] {
	return func(observe Observer[U], scheduler Scheduler, subscriber Subscriber) {
		state := seed
		observer := func(next T, err error, done bool) {
			switch {
			case !done:
				state = accumulator(state, next)
			case err != nil:
				var zero U
				observe(zero, err, true)
			default:
				Of(state)(observe, scheduler, subscriber)
			}
		}
		observable(observer, scheduler, subscriber)
	}
}

func (observable Observable[T]) Reduce(seed T, accumulator func(acc, next T) T) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		state := seed
		observer := func(next T, err error, done bool) {
			switch {
			case !done:
				state = accumulator(state, next)
			case err != nil:
				var zero T
				observe(zero, err, true)
			default:
				Of(state)(observe, scheduler, subscriber)
			}
		}
		observable(observer, scheduler, subscriber)
	}
}
