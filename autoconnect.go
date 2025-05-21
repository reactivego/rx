package rx

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrInvalidCount = errors.Join(Err, errors.New("invalid count"))

// AutoConnect returns an Observable that automatically connects to the Connectable source when a specified
// number of subscribers subscribe to it.
//
// When the specified number of subscribers (count) is reached, the Connectable source is connected,
// allowing it to start emitting items. The connection is shared among all subscribers.
// When all subscribers unsubscribe, the connection is terminated.
//
// If count is less than 1, it returns an Observable that emits an ErrInvalidCount error.
func (connectable Connectable[T]) AutoConnect(count int) Observable[T] {
	if count < 1 {
		return Throw[T](ErrInvalidCount)
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
				connectable.Connector(scheduler, source.subscription)
				source.Lock()
			}
		}
		source.Unlock()
	}
	return observable
}
