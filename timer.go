package rx

import "time"

func Timer[T Integer | Float](initialDelay time.Duration, intervals ...time.Duration) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var i T
		task := func(again func(due time.Duration)) {
			if subscriber.Subscribed() {
				if i == 0 || (i > 0 && len(intervals) > 0) {
					observe(i, nil, false)
				}
				if subscriber.Subscribed() {
					if len(intervals) > 0 {
						again(intervals[int(i)%len(intervals)])
					} else {
						if i == 0 {
							again(0)
						} else {
							var zero T
							observe(zero, nil, true)
						}
					}
				}
				i++
			}
		}
		runner := scheduler.ScheduleFutureRecursive(initialDelay, task)
		subscriber.OnUnsubscribe(runner.Cancel)
	}
}
