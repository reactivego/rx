package rx

// Take returns an Observable that emits only the first count values emitted by
// the source Observable. If the source emits fewer than count values then all
// of its values are emitted. After that, it completes, regardless if the source
// completes.
func Take[T any](n int) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
		if n <= 0 {
			return Empty[T]()
		}
		return Observable[T](func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
			taken := 0
			observable(func(next T, err error, done bool) {
				if !done {
					observe(next, nil, false)
					if taken++; taken == n {
						var zero T
						observe(zero, nil, true)
					}
				} else {
					observe(next, err, true)
				}
			}, scheduler, subscriber)
		}).AutoUnsubscribe()
	}
}

// Take returns an Observable that emits only the first count values emitted by
// the source Observable. If the source emits fewer than count values then all
// of its values are emitted. After that, it completes, regardless if the source
// completes.
func (observable Observable[T]) Take(n int) Observable[T] {
	return Take[T](n)(observable)
}
