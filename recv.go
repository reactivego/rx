package rx

func Recv[T any](ch <-chan T) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.ScheduleRecursive(func(again func()) {
			if !subscriber.Subscribed() {
				return
			}
			select {
			case next, ok := <-ch:
				if !subscriber.Subscribed() {
					return
				}
				if ok {
					err, ok := any(next).(error)
					if !ok {
						observe(next, nil, false)
						if !subscriber.Subscribed() {
							return
						}
						again()
					} else {
						var zero T
						observe(zero, err, true)
					}
				} else {
					var zero T
					observe(zero, nil, true)
				}
			case <-subscriber.Context().Done():
				return
			}
		})
		subscriber.OnUnsubscribe(func() {
			runner.Cancel()
		})
	}
}
