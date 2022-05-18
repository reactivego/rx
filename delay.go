package x

import (
	"sync"
	"time"
)

type emission[T any] struct {
	at   time.Time
	next T
	err  error
	done bool
}

func (observable Observable[T]) Delay(duration time.Duration) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var delay struct {
			sync.Mutex
			emissions []emission[T]
		}
		delayer := scheduler.ScheduleFutureRecursive(duration, func(again func(time.Duration)) {
			if subscriber.Subscribed() {
				delay.Lock()
				for _, entry := range delay.emissions {
					delay.Unlock()
					due := entry.at.Sub(scheduler.Now())
					if due > 0 {
						again(due)
						return
					}
					observe(entry.next, entry.err, entry.done)
					if entry.done || !subscriber.Subscribed() {
						return
					}
					delay.Lock()
					delay.emissions = delay.emissions[1:]
				}
				delay.Unlock()
				again(duration) // keep on rescheduling the emitter
			}
		})
		subscriber.OnUnsubscribe(delayer.Cancel)
		observer := func(next T, err error, done bool) {
			delay.Lock()
			delay.emissions = append(delay.emissions, emission[T]{scheduler.Now().Add(duration), next, err, done})
			delay.Unlock()
		}
		observable(observer, scheduler, subscriber)
	}
}
