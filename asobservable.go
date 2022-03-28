package observable

const TypecastFailed = Error("typecast failed")

func AsObservable[T any](observable Observable[any]) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		observer := func(next any, err error, done bool) {
			if !done {
				if nextT, ok := next.(T); ok {
					observe(nextT, err, done)
				} else {
					var zero T
					observe(zero, TypecastFailed, true)
				}
			} else {
				var zero T
				observe(zero, err, true)
			}
		}
		observable(observer, scheduler, subscriber)
	}
}

func (observable Observable[T]) AsObservable() Observable[any] {
	return func(observe Observer[any], scheduler Scheduler, subscriber Subscriber) {
		observer := func(next T, err error, done bool) {
			observe(next, err, done)
		}
		observable(observer, scheduler, subscriber)
	}
}
