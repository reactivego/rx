package x

func (observable Observable[T]) TakeWhile(condition func(next T) bool) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		observable(func(next T, err error, done bool) {
			if done || condition(next) {
				observe(next, err, done)
			} else {
				var zero T
				observe(zero, nil, true)
			}
		}, scheduler, subscriber)
	}
}
