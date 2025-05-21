package rx

import (
	"sync"
	"sync/atomic"
)

func MergeAll[T any](observable Observable[Observable[T]]) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var merge struct {
			sync.Mutex
			done  bool
			count atomic.Int32
		}
		merger := func(next T, err error, done bool) {
			merge.Lock()
			defer merge.Unlock()
			if !merge.done {
				switch {
				case !done:
					observe(next, nil, false)
				case err != nil:
					merge.done = true
					var zero T
					observe(zero, err, true)
				default:
					if merge.count.Add(-1) == 0 {
						var zero T
						observe(zero, nil, true)
					}
				}
			}
		}
		appender := func(next Observable[T], err error, done bool) {
			if !done {
				merge.count.Add(1)
				next.AutoUnsubscribe()(merger, scheduler, subscriber)
			} else {
				var zero T
				merger(zero, err, true)
			}
		}
		merge.count.Add(1)
		observable.AutoUnsubscribe()(appender, scheduler, subscriber)
	}
}
