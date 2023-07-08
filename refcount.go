package x

import (
	"sync"
	"sync/atomic"
)

// RefCount makes a Connectable[T] behave like an ordinary Observable[T].
// On first Subscribe it will call Connect on its Connectable[T] and when
// its last subscriber is Unsubscribed it will cancel the source connection by
// calling Unsubscribe on the subscription returned by the call to Connect.
func (connectable Connectable[T]) RefCount() Observable[T] {
	var source struct {
		sync.Mutex
		refcount     int32
		subscription *subscription
	}
	observable := func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		subscriber.OnUnsubscribe(func() {
			source.Lock()
			if atomic.AddInt32(&source.refcount, -1) == 0 {
				source.subscription.Unsubscribe()
			}
			source.Unlock()
		})
		connectable.Observable(observe, scheduler, subscriber)
		source.Lock()
		if atomic.AddInt32(&source.refcount, 1) == 1 {
			source.subscription = newSubscription(scheduler)
			source.Unlock()
			connectable.Connect(scheduler, source.subscription)
			source.Lock()
		}
		source.Unlock()
	}
	return observable
}
