package observable

func (observable Observable[T]) Skip(n int) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		i := 0
		observable(func(next T, err error, done bool) {
			if done || i >= n {
				observe(next, err, done)
			}
			i++
		}, scheduler, subscriber)
	}
}
