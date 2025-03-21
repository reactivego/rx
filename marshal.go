package rx

func (observable Observable[T]) Marshal(marshal func(any) ([]byte, error)) Observable[[]byte] {
	return func(observe Observer[[]byte], scheduler Scheduler, subscriber Subscriber) {
		observer := func(next T, err error, done bool) {
			if !done {
				data, err := marshal(next)
				observe(data, err, err != nil)
			} else {
				observe(nil, err, true)
			}
		}
		observable(observer, scheduler, subscriber)
	}
}
