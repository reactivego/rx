package x

func Create[T any](create Creator[T]) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		task := func(index int, again func(int)) {
			if subscriber.Subscribed() {
				next, err, done := create(index)
				if subscriber.Subscribed() {
					if !done {
						observe(next, nil, false)
						if subscriber.Subscribed() {
							again(index + 1)
						}
					} else {
						var zero T
						observe(zero, err, true)
					}
				}
			}
		}
		runner := scheduler.ScheduleLoop(0, task)
		subscriber.OnUnsubscribe(runner.Cancel)
	}
}
