package rx

import (
	"sync"
)

func ConcatAll[T any](observable Observable[Observable[T]]) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var concat struct {
			sync.Mutex
			chain []func()
		}
		concatenator := func(next T, err error, done bool) {
			concat.Lock()
			defer concat.Unlock()
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if len(concat.chain) > 0 {
					link := concat.chain[0]
					concat.chain = concat.chain[1:]
					link()
				} else {
					concat.chain = nil
				}
			}
		}
		link := func(observable Observable[T]) func() {
			return func() {
				if observable != nil {
					observable.AutoUnsubscribe()(concatenator, scheduler, subscriber)
				} else {
					Empty[T]()(observe, scheduler, subscriber)
				}
			}
		}
		appender := func(next Observable[T], err error, done bool) {
			if !done || err == nil {
				concat.Lock()
				active := (concat.chain != nil)
				concat.chain = append(concat.chain, link(next))
				concat.Unlock()
				if !active {
					var zero T
					concatenator(zero, nil, true)
				}
			} else {
				var zero T
				concatenator(zero, err, true)
			}
		}
		observable.AutoUnsubscribe()(appender, scheduler, subscriber)
	}
}
