package rx

import (
	"sync"
	"sync/atomic"
)

//jig:template Observable<Foo> Do

// Do calls a function for each next value passing through the observable.
func (o ObservableFoo) Do(f func(next foo)) ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
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
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
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
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
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

//jig:template Observable<Foo> Finally

// Finally applies a function for any error or completion on the stream.
// This doesn't expose whether this was an error or a completion.
func (o ObservableFoo) Finally(f func()) ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
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
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		// Subscribe scope
		observer := func(next foo, err error, done bool) {
			// Observe scope
			observe(next, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable Repeat

// Repeat creates an Observable that emits a sequence of items repeatedly.
func (o Observable) Repeat(count int) Observable {
	if count == 0 {
		return Empty()
	}
	observable := func(observe Observer, scheduler Scheduler, subscriber Subscriber) {
		var repeated int
		var observer Observer
		observer = func(next interface{}, err error, done bool) {
			if !done || err != nil {
				observe(next, err, done)
			} else {
				repeated++
				if repeated < count {
					o(observer, scheduler, subscriber)
				} else {
					observe(nil, nil, true)
				}
			}
		}
		o(observer, scheduler, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Repeat
//jig:needs Observable Repeat

// Repeat creates an ObservableFoo that emits a sequence of items repeatedly.
func (o ObservableFoo) Repeat(count int) ObservableFoo {
	return o.AsObservable().Repeat(count).AsObservableFoo()
}

//jig:template Observable<Foo> Serialize

// Serialize forces an ObservableFoo to make serialized calls and to be
// well-behaved.
func (o ObservableFoo) Serialize() ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
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
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
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
					var zeroFoo foo
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
