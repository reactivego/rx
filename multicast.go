package x

import "sync"

// Multicast returns both an Observer and and Observable. The returned Observer is
// used to send items into the Multicast. The returned Observable is used to subscribe
// to the Multicast. The Multicast multicasts items send through the Observer to every
// Subscriber of the Observable.
//
//	size  size of the item buffer, number of items kept to replay to a new Subscriber.
//
// Backpressure handling depends on the sign of the size argument. For positive size the
// multicast will block when one of the subscribers lets the buffer fill up. For negative
// size the multicast will drop items on the blocking subscriber, allowing the others to
// keep on receiving values. For hot observables dropping is preferred.
func Multicast[T any](size int) (Observer[T], Observable[T]) {
	var multicast struct {
		sync.Mutex
		channels []chan any
	}
	drop := (size < 0)
	if size < 0 {
		size = -size
	}
	observer := func(next T, err error, done bool) {
		multicast.Lock()
		defer multicast.Unlock()
		for _, c := range multicast.channels {
			if c != nil {
				switch {
				case !done:
					if drop {
						select {
						case c <- next:
						default:
							// dropping
						}
					} else {
						c <- next
					}
				case err != nil:
					if drop {
						select {
						case c <- err:
						default:
							// dropping
						}
					} else {
						c <- err
					}
					close(c)
				default:
					close(c)
				}
			}
		}
	}
	observable := func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		multicast.Lock()
		defer multicast.Unlock()
		channel := make(chan any, size)
		remove := func() bool {
			for i, v := range multicast.channels {
				if v == channel {
					copy(multicast.channels[i:], multicast.channels[i+1:])
					size := len(multicast.channels) - 1
					multicast.channels[size] = nil
					multicast.channels = multicast.channels[:size]
					return true
				}
			}
			return false
		}
		multicast.channels = append(multicast.channels, channel)
		observer := func(next any, err error, done bool) {
			switch {
			case !done:
				switch n := next.(type) {
				case error:
					if remove() {
						var zero T
						observe(zero, n, true)
					}
				case T:
					observe(n, nil, false)
				}
			case err != nil:
				if remove() {
					var zero T
					observe(zero, err, true)
				}
			default:
				if remove() {
					var zero T
					observe(zero, nil, true)
				}
			}
		}
		Recv(channel)(observer, scheduler, subscriber)
		subscriber.OnUnsubscribe(func() { remove() })
	}
	return observer, observable
}
