package x

func (observable Observable[T]) Do(f func(next T)) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		observer := func(next T, err error, done bool) {
			if !done {
				f(next)
			}
			observe(next, err, done)
		}
		observable(observer, scheduler, subscriber)
	}
}
