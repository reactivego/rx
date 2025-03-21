package rx

func From[T any](slice ...T) Observable[T] {
	if len(slice) == 0 {
		return Empty[T]()
	}
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		task := func(index int, again func(next int)) {
			if subscriber.Subscribed() {
				if index < len(slice) {
					observe(slice[index], nil, false)
					if subscriber.Subscribed() {
						again(index + 1)
					}
				} else {
					var zero T
					observe(zero, nil, true)
				}
			}
		}
		runner := scheduler.ScheduleLoop(0, task)
		subscriber.OnUnsubscribe(runner.Cancel)
	}
}
