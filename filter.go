package observable

func (observable Observable[T]) Filter(predicate func(T) bool) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		observable(func(next T, err error, done bool) {
			if done || predicate(next) {
				observe(next, err, done)
			}
		}, scheduler, subscriber)
	}
}
