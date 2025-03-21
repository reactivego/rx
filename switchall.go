package rx

import (
	"sync"
	"sync/atomic"
)

func SwitchAll[T any](observable Observable[Observable[T]]) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var switches struct {
			sync.Mutex
			done       bool
			count      int32
			subscriber Subscriber
		}
		observer := func(next T, err error, done bool) {
			switches.Lock()
			defer switches.Unlock()
			if !switches.done {
				switch {
				case !done:
					observe(next, nil, false)
				case err != nil:
					switches.done = true
					var zero T
					observe(zero, err, true)
				default:
					if atomic.AddInt32(&switches.count, -1) == 0 {
						var zero T
						observe(zero, nil, true)
					}
				}
			}
		}
		switcher := func(next Observable[T], err error, done bool) {
			if !done {
				switch {
				case atomic.CompareAndSwapInt32(&switches.count, 2, 3):
					switches.Lock()
					defer switches.Unlock()
					switches.subscriber.Unsubscribe()
					atomic.StoreInt32(&switches.count, 2)
					fallthrough
				case atomic.CompareAndSwapInt32(&switches.count, 1, 2):
					switches.subscriber = subscriber.Add()
					next.AutoUnsubscribe()(observer, scheduler, switches.subscriber)
				}
			} else {
				var zero T
				observer(zero, err, true)
			}
		}
		switches.count += 1
		observable.AutoUnsubscribe()(switcher, scheduler, subscriber)
	}
}
