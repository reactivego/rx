package observable

func (observable Observable[T]) Count() Observable[int] {
	return func(observe Observer[int], scheduler Scheduler, subscriber Subscriber) {
		var count int
		observer := func(next T, err error, done bool) {
			if !done {
				count++
			} else {
				Of(count)(observe, scheduler, subscriber)
			}
		}
		observable(observer, scheduler, subscriber)
	}
}
