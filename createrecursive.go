package observable

func CreateRecursive[T any](create Creator[T]) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		task := func(again func()) {
			if subscriber.Subscribed() {
				next, err, done := create()
				if subscriber.Subscribed() {
					if !done {
						observe(next, nil, false)
						if subscriber.Subscribed() {
							again()
						}
					} else {
						var zero T
						observe(zero, err, true)
					}
				}
			}
		}
		runner := scheduler.ScheduleRecursive(task)
		subscriber.OnUnsubscribe(runner.Cancel)
	}
}
