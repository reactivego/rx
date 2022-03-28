package observable

func (observable Observable[T]) Catch(other Observable[T]) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		observer := func(next T, err error, done bool) {
			if !done || err == nil {
				observe(next, err, done)
			} else {
				other(observe, scheduler, subscriber)
			}
		}
		observable(observer, scheduler, subscriber)
	}
}
