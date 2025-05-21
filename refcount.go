package rx

import (
	"sync"
	"sync/atomic"
)

// RefCount converts a Connectable Observable into a standard Observable that automatically
// connects when the first subscriber subscribes and disconnects when the last subscriber
// unsubscribes.
//
// When the first subscriber subscribes to the resulting Observable, it automatically calls
// Connect() on the source Connectable Observable. The connection is shared among all
// subscribers. When the last subscriber unsubscribes, the connection is automatically
// closed.
//
// This is useful for efficiently sharing expensive resources (like network connections)
// among multiple subscribers.
func (connectable Connectable[T]) RefCount() Observable[T] {
	var source struct {
		sync.Mutex
		refcount     int32
		subscription *subscription
	}
	observable := func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		source.Lock()
		if atomic.AddInt32(&source.refcount, 1) == 1 {
			source.subscription = newSubscription(scheduler)
			source.Unlock()
			connectable.Connector(scheduler, source.subscription)
			source.Lock()
		}
		source.Unlock()
		subscriber.OnUnsubscribe(func() {
			source.Lock()
			if atomic.AddInt32(&source.refcount, -1) == 0 {
				source.subscription.Unsubscribe()
			}
			source.Unlock()
		})
		connectable.Observable(observe, scheduler, subscriber)
	}
	return observable
}
