package rx

func (observable Observable[T]) ConcatWith(others ...Observable[T]) Observable[T] {
	if len(others) == 0 {
		return observable
	}
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var (
			observables = append([]Observable[T]{}, others...)
			observer    Observer[T]
		)
		observer = func(next T, err error, done bool) {
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if len(observables) == 0 {
					var zero T
					observe(zero, nil, true)
				} else {
					o := observables[0]
					observables = observables[1:]
					o(observer, scheduler, subscriber)
				}
			}
		}
		observable(observer, scheduler, subscriber)
	}
}
