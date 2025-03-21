package rx

func Empty[T any]() Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		task := func() {
			if subscriber.Subscribed() {
				var zero T
				observe(zero, nil, true)
			}
		}
		runner := scheduler.Schedule(task)
		subscriber.OnUnsubscribe(runner.Cancel)
	}
}
