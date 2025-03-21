package rx

import (
	"sync"
	"sync/atomic"
)

func ExhaustAll[T any](observable Observable[Observable[T]]) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var observers struct {
			sync.Mutex
			done  bool
			count int32
		}
		observer := func(next T, err error, done bool) {
			observers.Lock()
			defer observers.Unlock()
			if !observers.done {
				switch {
				case !done:
					observe(next, nil, false)
				case err != nil:
					observers.done = true
					var zero T
					observe(zero, err, true)
				default:
					if atomic.AddInt32(&observers.count, -1) == 0 {
						var zero T
						observe(zero, nil, true)
					}
				}
			}
		}
		exhauster := func(next Observable[T], err error, done bool) {
			if !done {
				if atomic.CompareAndSwapInt32(&observers.count, 1, 2) {
					next.AutoUnsubscribe()(observer, scheduler, subscriber)
				}
			} else {
				var zero T
				observer(zero, err, true)
			}
		}
		observers.count += 1
		observable.AutoUnsubscribe()(exhauster, scheduler, subscriber)
	}
}
