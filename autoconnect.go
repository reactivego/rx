package rx

import (
	"sync"
	"sync/atomic"
)

const InvalidCount = Error("invalid count")

func (connectable Connectable[T]) AutoConnect(count int) Observable[T] {
	if count < 1 {
		return Throw[T](InvalidCount)
	}
	var source struct {
		sync.Mutex
		refcount     int32
		subscription *subscription
	}
	observable := func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		subscriber.OnUnsubscribe(func() {
			source.Lock()
			if atomic.AddInt32(&source.refcount, -1) == 0 {
				if source.subscription != nil {
					source.subscription.Unsubscribe()
				}
			}
			source.Unlock()
		})
		connectable.Observable(observe, scheduler, subscriber)
		source.Lock()
		if atomic.AddInt32(&source.refcount, 1) == int32(count) {
			if source.subscription == nil || source.subscription.err != nil {
				source.subscription = newSubscription(scheduler)
				source.Unlock()
				connectable.Connect(scheduler, source.subscription)
				source.Lock()
			}
		}
		source.Unlock()
	}
	return observable
}
