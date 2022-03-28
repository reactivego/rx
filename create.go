package observable

func Create[T any](create func() (T, error)) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		task := func() {
			if subscriber.Subscribed() {
				if next, err := create(); err == nil {
					Of(next)(observe, scheduler, subscriber)
				} else {
					Throw[T](err)(observe, scheduler, subscriber)
				}
			}
		}
		runner := scheduler.Schedule(task)
		subscriber.OnUnsubscribe(runner.Cancel)
	}
}
