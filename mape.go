package x

func MapE[T, U any](observable Observable[T], project func(T) (U, error)) Observable[U] {
	return func(observe Observer[U], scheduler Scheduler, subscriber Subscriber) {
		observable(func(next T, err error, done bool) {
			var mapped U
			if !done {
				mapped, err = project(next)
			}
			observe(mapped, err, done || err != nil)
		}, scheduler, subscriber)
	}
}
