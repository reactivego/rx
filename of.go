package rx

func Of[T any](value T) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		task := func(index int, again func(next int)) {
			if subscriber.Subscribed() {
				if index == 0 {
					observe(value, nil, false)
					if subscriber.Subscribed() {
						again(1)
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
