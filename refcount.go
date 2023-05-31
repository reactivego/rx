package x

import (
	"sync"
	"sync/atomic"
)

// RefCount makes a ConnectableObservable[T] behave like an ordinary Observable[T].
// On first Subscribe it will call Connect on its ConnectableObservable[T] and when
// its last subscriber is Unsubscribed it will cancel the source connection by
// calling Unsubscribe on the subscription returned by the call to Connect.
func (o ConnectableObservable[T]) RefCount() Observable[T] {
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
		o.Observable(observe, scheduler, subscriber)
		source.Lock()
		if atomic.AddInt32(&source.refcount, 1) == 1 {
			source.subscription = newSubscription(scheduler)
			source.Unlock()
			o.Connectable(scheduler, source.subscription)
			source.Lock()
		}
		source.Unlock()
	}
	return observable
}
