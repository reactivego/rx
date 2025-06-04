package rx

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
		err      error
		done     bool
	}
	drop := (size < 0)
	if size < 0 {
		size = -size
	}
	observer := func(next T, err error, done bool) {
		multicast.Lock()
		defer multicast.Unlock()
		if !multicast.done {
			for _, ch := range multicast.channels {
				switch {
				case !done:
					if drop {
						select {
						case ch <- next:
						default:
							// dropping
						}
					} else {
						ch <- next
					}
				case err != nil:
					if drop {
						select {
						case ch <- err:
						default:
							// dropping
						}
					} else {
						ch <- err
					}
					close(ch)
				default:
					close(ch)
				}
			}
			multicast.err = err
			multicast.done = done
		}
	}
	observable := func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		addChannel := func() <-chan any {
			multicast.Lock()
			defer multicast.Unlock()
			var ch chan any
			if !multicast.done {
				ch = make(chan any, size)
				multicast.channels = append(multicast.channels, ch)
			} else if multicast.err != nil {
				ch = make(chan any, 1)
				ch <- multicast.err
				close(ch)
			} else {
				ch = make(chan any)
				close(ch)
			}
			return ch
		}
		removeChannel := func(ch <-chan any) func() {
			return func() {
				multicast.Lock()
				defer multicast.Unlock()
				if !multicast.done && ch != nil {
					channels := multicast.channels[:0]
					for _, c := range multicast.channels {
						if ch != c {
							channels = append(channels, c)
						}
					}
					multicast.channels = channels
				}
			}
		}
		ch := addChannel()
		Recv(ch)(observe.AsObserver(), scheduler, subscriber)
		subscriber.OnUnsubscribe(removeChannel(ch))
	}
	return observer, observable
}
