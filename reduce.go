package rx

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
