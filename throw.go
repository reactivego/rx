package observable

func Throw[T any](err error) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		task := func() {
			if subscriber.Subscribed() {
				var zero T
				observe(zero, err, true)
			}
		}
		runner := scheduler.Schedule(task)
		subscriber.OnUnsubscribe(runner.Cancel)
	}
}

func Throww[T any](err error) Observable[T] {
	return Throw[T](Errorw(err))
}
