package x

import (
	"math"
	"time"
)

func (observable Observable[T]) Retry(count ...int) Observable[T] {
	if len(count) == 0 || count[0] <= 0 {
		count = []int{math.MaxInt}
	}
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var retry struct {
			count       int
			observer    Observer[T]
			subscriber  Subscriber
			resubscribe func()
		}
		retry.count = count[0]
		retry.observer = func(next T, err error, done bool) {
			switch {
			case !done:
				observe(next, nil, false)
			case err != nil && retry.count > 0:
				retry.count--
				retry.subscriber.Unsubscribe()
				scheduler.ScheduleFuture(1*time.Millisecond, retry.resubscribe)
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
