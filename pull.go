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

func Pull2[K, V any](seq iter.Seq2[K, V]) Observable[Tuple2[K, V]] {
	return func(observe Observer[Tuple2[K, V]], scheduler Scheduler, subscriber Subscriber) {
		next, stop := iter.Pull2(seq)
		task := func(again func()) {
			if subscriber.Subscribed() {
				k, v, valid := next()
				if subscriber.Subscribed() {
					if valid {
						observe(Tuple2[K, V]{k, v}, nil, false)
						if subscriber.Subscribed() {
							again()
						}
					} else {
						var zero Tuple2[K, V]
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
