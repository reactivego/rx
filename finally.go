package x

func (observable Observable[T]) Finally(f func(error)) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		observer := func(next T, err error, done bool) {
			if done {
				f(err)
			}
			observe(next, err, done)
		}
		observable(observer, scheduler, subscriber)
	}
}
