package rx

import "sync"

func (observable Observable[T]) MergeWith(others ...Observable[T]) Observable[T] {
	if len(others) == 0 {
		return observable
	}
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var merge struct {
			sync.Mutex
			done  bool
			count int
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
					if merge.count--; merge.count == 0 {
						var zero T
						observe(zero, nil, true)
					}
				}
			}
		}
		merge.count = 1 + len(others)
		observable.AutoUnsubscribe()(merger, scheduler, subscriber)
		for _, other := range others {
			if subscriber.Subscribed() {
				other.AutoUnsubscribe()(merger, scheduler, subscriber)
			}
		}
	}
}
