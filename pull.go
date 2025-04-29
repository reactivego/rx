package rx

import "iter"

func Pull[T any](seq iter.Seq[T]) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		next, stop := iter.Pull(seq)
		task := func(again func()) {
			if subscriber.Subscribed() {
				next, valid := next()
				if subscriber.Subscribed() {
					if valid {
						observe(next, nil, false)
						if subscriber.Subscribed() {
							again()
						}
					} else {
						observe(next, nil, true)
					}
				}
			}
		}
		runner := scheduler.ScheduleRecursive(task)
		subscriber.OnUnsubscribe(runner.Cancel)
		subscriber.OnUnsubscribe(stop)
	}
}

func Pull2[T, U any](seq iter.Seq2[T, U]) Observable[Tuple2[T, U]] {
	return func(observe Observer[Tuple2[T, U]], scheduler Scheduler, subscriber Subscriber) {
		next, stop := iter.Pull2(seq)
		task := func(again func()) {
			if subscriber.Subscribed() {
				first, second, valid := next()
				if subscriber.Subscribed() {
					if valid {
						observe(Tuple2[T, U]{first, second}, nil, false)
						if subscriber.Subscribed() {
							again()
						}
					} else {
						var zero Tuple2[T, U]
						observe(zero, nil, true)
					}
				}
			}
		}
		runner := scheduler.ScheduleRecursive(task)
		subscriber.OnUnsubscribe(runner.Cancel)
		subscriber.OnUnsubscribe(stop)
	}
}
