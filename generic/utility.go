package rx

import (
	"sync"
	"sync/atomic"
	"time"
)

//jig:template Observable<Foo> Do

// Do calls a function for each next value passing through the observable.
func (o ObservableFoo) Do(f func(next foo)) ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			if !done {
				f(next)
			}
			observe(next, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> DoOnError

// DoOnError calls a function for any error on the stream.
func (o ObservableFoo) DoOnError(f func(err error)) ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			if err != nil {
				f(err)
			}
			observe(next, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> DoOnComplete

// DoOnComplete calls a function when the stream completes.
func (o ObservableFoo) DoOnComplete(f func()) ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			if err == nil && done {
				f()
			}
			observe(next, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Delay

// Delay shifts the emission from an Observable forward in time by a particular amount of time.
func (o ObservableFoo) Delay(duration time.Duration) ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		firstTime := true
		observer := func(next foo, err error, done bool) {
			if firstTime {
				firstTime = false
				time.Sleep(duration)
			}
			observe(next, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Finally

// Finally applies a function for any error or completion on the stream.
// This doesn't expose whether this was an error or a completion.
func (o ObservableFoo) Finally(f func()) ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			if done {
				f()
			}
			observe(next, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Passthrough

// Passthrough just passes through all output from the ObservableFoo.
func (o ObservableFoo) Passthrough() ObservableFoo {
	// Operator scope
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		// Subscribe scope
		observer := func(next foo, err error, done bool) {
			// Observe scope
			observe(next, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Serialize

// Serialize forces an ObservableFoo to make serialized calls and to be
// well-behaved.
func (o ObservableFoo) Serialize() ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var observer struct {
			sync.Mutex
			done bool
		}
		serializer := func(next foo, err error, done bool) {
			observer.Lock()
			defer observer.Unlock()
			if !observer.done {
				observer.done = done
				observe(next, err, done)
			}
		}
		o(serializer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable Timeout
//jig:needs RxError, Observable Serialize

// ErrTimeout is delivered to an observer if the stream times out.
const ErrTimeout = RxError("timeout")

// Timeout mirrors the source Observable, but issues an error notification if a
// particular period of time elapses without any emitted items.
// Timeout schedules tasks on the scheduler passed to this 
func (o Observable) Timeout(timeout time.Duration) Observable {
	observable := Observable(func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		if subscriber.Canceled() {
			return
		}
		var last struct {
			sync.Mutex
			at time.Time
			done bool
		}
		last.at = subscribeOn.Now()
		timer := func(self func(time.Duration)){
			last.Lock()
			defer last.Unlock()
			if last.done || subscriber.Canceled() {
				return
			}
			deadline := last.at.Add(timeout)
			now := subscribeOn.Now()
			if now.Before(deadline) {
				self(deadline.Sub(now))
				return
			}
			last.done = true
			observe(zero, ErrTimeout, true)
		}
		runner := subscribeOn.ScheduleFutureRecursive(timeout, timer)
		subscriber.OnUnsubscribe(runner.Cancel)
		observer := func(next interface{}, err error, done bool) {
			last.Lock()
			defer last.Unlock()
			if last.done || subscriber.Canceled() {
				return
			}
			now := subscribeOn.Now()
			deadline := last.at.Add(timeout) 
			if !now.Before(deadline) {
				return
			}
			last.done = done
			last.at = now
			observe(next, err, done)
		}
		o(observer, subscribeOn, subscriber)
	})
	return observable.Serialize()
}

//jig:template Observable<Foo> Timeout
//jig:needs Observable Timeout

// Timeout mirrors the source ObservableFoo, but issues an error notification if
// a particular period of time elapses without any emitted items.
//
// This observer starts a goroutine for every subscription to monitor the
// timeout deadline. It is guaranteed that calls to the observer for this
// subscription will never be called concurrently. It is however almost certain
// that any timeout error will be delivered on a goroutine other than the one
// delivering the next values.
func (o ObservableFoo) Timeout(timeout time.Duration) ObservableFoo {
	return o.AsObservable().Timeout(timeout).AsObservableFoo()
}

//jig:template ErrObservableContractViolation
//jig:needs RxError

const ErrObservableContractViolationConcurrentNotifications = RxError("observable contract violation: concurrent notifications")
const ErrObservableContractViolationNextAfterTermination = RxError("observable contract violation: next after termination")
const ErrObservableContractViolationErrorAfterTermination = RxError("observable contract violation: error after termination")
const ErrObservableContractViolationCompleteAfterTermination = RxError("observable contract violation: complete after termination")

//jig:template Observable<Foo> Validated
//jig:needs ErrObservableContractViolation

// Validated will check for violations of the observable contract. More specific
// it will detect concurrent notifications from the observable and it will
// detect notifications sent after an error or complete notification was already
// processed. A violation is reported via the ObservableFoo when the observable
// is not already terminated (we don't want to break the contract ourselves).
// A violation will always be reported via the onViolation callback. If nil is
// passed as the callback, the operation will use panic to report the violation.
func (o ObservableFoo) Validated(onViolation func(err error)) ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		const (
			operational int32 = iota
			terminated
			violation
		)
		nextstate := func(done bool) int32 {
			if done {
				return terminated
			} else {
				return operational
			}
		}
		if onViolation == nil {
			onViolation = func(err error) { panic(err) }
		}
		concurrent := int32(0)
		state := operational
		var mu sync.Mutex
		observer := func(next foo, err error, done bool) {
			if atomic.AddInt32(&concurrent, 1) > 1 {
				mu.Lock()
				if atomic.CompareAndSwapInt32(&state, operational, violation) {
					err = ErrObservableContractViolationConcurrentNotifications
					observe(zeroFoo, err, true)
					onViolation(err)
				} else if atomic.CompareAndSwapInt32(&state, terminated, violation) {
					onViolation(ErrObservableContractViolationConcurrentNotifications)
				}
				mu.Unlock()
			} else {
				mu.Lock()
				if atomic.CompareAndSwapInt32(&state, terminated, violation) {
					if !done {
						onViolation(ErrObservableContractViolationNextAfterTermination)
					} else {
						if err != nil {
							onViolation(ErrObservableContractViolationErrorAfterTermination)
						} else {
							onViolation(ErrObservableContractViolationCompleteAfterTermination)
						}
					}
				} else if atomic.CompareAndSwapInt32(&state, operational, nextstate(done)) {
					observe(next, err, done)
				}
				mu.Unlock()
			}
			atomic.AddInt32(&concurrent, -1)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}
