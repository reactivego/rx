package observable

import (
	"sync"
	"time"
)

// SampleTime emits the most recent item emitted by an Observable within periodic time intervals.
func (observable Observable[T]) SampleTime(window time.Duration) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var sample struct {
			sync.Mutex
			at   time.Time
			next T
			done bool
		}
		sampler := scheduler.ScheduleFutureRecursive(window, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				sample.Lock()
				if !sample.done {
					begin := scheduler.Now().Add(-window)
					if !sample.at.Before(begin) {
						observe(sample.next, nil, false)
					}
					if subscriber.Subscribed() {
						self(window)
					}
				}
				sample.Unlock()
			}
		})
		subscriber.OnUnsubscribe(sampler.Cancel)
		observer := func(next T, err error, done bool) {
			if subscriber.Subscribed() {
				sample.Lock()
				sample.at = scheduler.Now()
				sample.next = next
				sample.done = done
				sample.Unlock()
				if done {
					var zero T
					observe(zero, err, true)
				}
			}
		}
		observable(observer, scheduler, subscriber)
	}
}
