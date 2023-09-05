package x

import (
	"math"
	"time"
)

func (observable Observable[T]) RetryTime(backoff func(int) time.Duration, limit ...int) Observable[T] {
	if len(limit) == 0 || limit[0] <= 0 {
		limit = []int{math.MaxInt}
	}
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var retry struct {
			observer    Observer[T]
			count       int
			subscriber  Subscriber
			resubscribe func()
		}
		retry.observer = func(next T, err error, done bool) {
			switch {
			case !done:
				observe(next, nil, false)
				retry.count = 0
			case err != nil && backoff != nil && retry.count < limit[0]:
				retry.subscriber.Unsubscribe()
				scheduler.ScheduleFuture(backoff(retry.count), retry.resubscribe)
				retry.count++
			default:
				observe(next, err, true)
			}
		}
		retry.resubscribe = func() {
			if subscriber.Subscribed() {
				retry.subscriber = subscriber.Add()
				observable(retry.observer, scheduler, retry.subscriber)
			}
		}
		retry.resubscribe()
	}
}
