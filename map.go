package observable

func Map[T, U any](observable Observable[T], project func(T) U) Observable[U] {
	return func(observe Observer[U], scheduler Scheduler, subscriber Subscriber) {
		observable(func(next T, err error, done bool) {
			var mapped U
			if !done {
				mapped = project(next)
			}
			observe(mapped, err, done)
		}, scheduler, subscriber)
	}
}
