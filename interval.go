package x

import (
	"time"

	"golang.org/x/exp/constraints"
)

func Interval[T constraints.Integer | constraints.Float](interval time.Duration) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var i T
		task := func(again func(due time.Duration)) {
			if subscriber.Subscribed() {
				observe(i, nil, false)
				i++
				if subscriber.Subscribed() {
					again(interval)
				}
			}
		}
		runner := scheduler.ScheduleFutureRecursive(interval, task)
		subscriber.OnUnsubscribe(runner.Cancel)
	}
}
