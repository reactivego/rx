package rx

import "time"

//jig:template Make<Foo>Func

// MakeFooFunc is the signature of a function that can be passed to MakeFoo
// to implement an ObservableFoo.
type MakeFooFunc func(Next func(foo), Error func(error), Complete func())

//jig:template Make<Foo>
//jig:needs Observable<Foo>, Make<Foo>Func

// MakeFoo provides a way of creating an ObservableFoo from scratch by
// calling observer methods programmatically. A make function conforming to the
// MakeFooFunc signature will be called by MakeFoo provining a Next, Error and
// Complete function that can be called by the code that implements the
// Observable.
func MakeFoo(make MakeFooFunc) ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		done := false
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Canceled() {
				return
			}
			next := func(n foo) {
				if subscriber.Subscribed() {
					observe(n, nil, false)
				}
			}
			err := func(e error) {
				done = true
				if subscriber.Subscribed() {
					observe(zeroFoo, e, true)
				}
			}
			complete := func() {
				done = true
				if subscriber.Subscribed() {
					observe(zeroFoo, nil, true)
				}
			}
			make(next, err, complete)
			if !done && subscriber.Subscribed() {
				self()
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template MakeTimed<Foo>Func

// MakeTimedFooFunc is the signature of a function that can be passed to
// MakeTimedFoo to implement an ObservableFoo.
type MakeTimedFooFunc func(Next func(foo), Error func(error), Complete func()) time.Duration

//jig:template MakeTimed<Foo>
//jig:needs Observable<Foo>, MakeTimed<Foo>Func

// MakeTimedFoo provides a way of creating an ObservableFoo from scratch by
// calling observer methods programmatically. A make function conforming to the
// MakeFooFunc signature will be called by MakeFoo provining a Next, Error and
// Complete function that can be called by the code that implements the
// Observable. The timeout passed in determines the time between calling the
// make function. The time.Duration returned by MakeTimedFooFunc determines when
// to reschedule the next iteration.
func MakeTimedFoo(timeout time.Duration, make MakeTimedFooFunc) ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		done := false
		runner := scheduler.ScheduleFutureRecursive(timeout, func(self func(time.Duration)) {
			if subscriber.Canceled() {
				return
			}
			next := func(n foo) {
				if subscriber.Subscribed() {
					observe(n, nil, false)
				}
			}
			error := func(e error) {
				done = true
				if subscriber.Subscribed() {
					observe(zeroFoo, e, true)
				}
			}
			complete := func() {
				done = true
				if subscriber.Subscribed() {
					observe(zeroFoo, nil, true)
				}
			}
			timeout = make(next, error, complete)
			if !done && subscriber.Subscribed() {
				self(timeout)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template <Foo>Observer
//jig:needs <Foo>ObserveFuncMethods

// FooObserver is the interface used with CreateFoo when implementing a custom
// observable.
type FooObserver interface {
	// Next emits the next foo value.
	Next(foo)
	// Error signals an error condition.
	Error(error)
	// Complete signals that no more data is to be expected.
	Complete()
	// Subscribed returns true when the subscription is currently valid.
	Subscribed() bool
}

//jig:template Create<Foo>
//jig:needs Observable<Foo>, <Foo>Observer

// CreateFoo creates an Observable from scratch by calling observer methods
// programmatically.
func CreateFoo(f func(FooObserver)) ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if !subscriber.Subscribed() {
				return
			}
			observer := func(next foo, err error, done bool) {
				if subscriber.Subscribed() {
					observe(next, err, done)
				}
			}
			type ObserverSubscriber struct {
				FooObserveFunc
				Subscriber
			}
			f(&ObserverSubscriber{observer, subscriber})
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template Defer<Foo>
//jig:needs Observable<Foo>

// DeferFoo does not create the ObservableFoo until the observer subscribes,
// and creates a fresh ObservableFoo for each observer.
func DeferFoo(factory func() ObservableFoo) ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		factory()(observe, scheduler, subscriber)
	}
	return observable
}

//jig:template Empty<Foo>
//jig:needs Observable<Foo>

// EmptyFoo creates an Observable that emits no items but terminates normally.
func EmptyFoo() ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Subscribed() {
				observe(zeroFoo, nil, true)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template Error<Foo>
//jig:needs Observable<Foo>

// ErrorFoo creates an Observable that emits no items and terminates with an
// error.
func ErrorFoo(err error) ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Subscribed() {
				observe(zeroFoo, err, true)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template FromChan<Foo>
//jig:needs Observable<Foo>

// FromChanFoo creates an ObservableFoo from a Go channel of foo values.
// It's not possible for the code feeding into the channel to send an error.
// The feeding code can send zero or more foo items and then closing the
// channel will be seen as completion.
func FromChanFoo(ch <-chan foo) ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Canceled() {
				return
			}
			next, ok := <-ch
			if subscriber.Canceled() {
				return
			}
			if ok {
				observe(next, nil, false)
				if subscriber.Subscribed() {
					self()
				}
			} else {
				observe(zeroFoo, nil, true)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template FromChan
//jig:needs Observable

// FromChan creates an Observable from a Go channel of interface{}
// values. This allows the code feeding into the channel to send either an error
// or the next value. The feeding code can send zero or more items and then
// closing the channel will be seen as completion. When the feeding code sends
// an error into the channel, it should close the channel immediately to
// indicate termination with error.
func FromChan(ch <-chan interface{}) Observable {
	observable := func(observe ObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Canceled() {
				return
			}
			next, ok := <-ch
			if subscriber.Canceled() {
				return
			}
			if ok {
				err, ok := next.(error)
				if !ok {
					observe(next, nil, false)
					if subscriber.Subscribed() {
						self()
					}
				} else {
					observe(zero, err, true)
				}
			} else {
				observe(zero, nil, true)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template From<Foo>
//jig:needs Observable<Foo>

// FromFoo creates an ObservableFoo from multiple foo values passed in.
func FromFoo(slice ...foo) ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		i := 0
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Subscribed() {
				if i < len(slice) {
					observe(slice[i], nil, false)
					if subscriber.Subscribed() {
						i++
						self()
					}
				} else {
					observe(zeroFoo, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template Interval
//jig:needs ObservableInt

// Interval creates an ObservableInt that emits a sequence of integers spaced
// by a particular time interval. First integer is emitted after the first time
// interval expires.
func Interval(interval time.Duration) ObservableInt {
	observable := func(observe IntObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		i := 0
		runner := scheduler.ScheduleFutureRecursive(interval, func(self func(time.Duration)) {
			if subscriber.Canceled() {
				return
			}
			observe(i, nil, false)
			if subscriber.Canceled() {
				return
			}
			i++
			self(interval)
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template Just<Foo>
//jig:needs Observable<Foo>

// JustFoo creates an ObservableFoo that emits a particular item.
func JustFoo(element foo) ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		done := false
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Subscribed() {
				if !done {
					observe(element, nil, false)
					if subscriber.Subscribed() {
						done = true
						self()
					}
				} else {
					observe(zeroFoo, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template Never<Foo>
//jig:needs Observable<Foo>

// NeverFoo creates an ObservableFoo that emits no items and does't terminate.
func NeverFoo() ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
	}
	return observable
}

//jig:template Of<Foo>
//jig:needs Observable<Foo>

// OfFoo emits a variable amount of values in a sequence and then emits a
// complete notification.
func OfFoo(slice ...foo) ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		i := 0
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Subscribed() {
				if i < len(slice) {
					observe(slice[i], nil, false)
					if subscriber.Subscribed() {
						i++
						self()
					}
				} else {
					observe(zeroFoo, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template Range
//jig:needs ObservableInt

// Range creates an ObservableInt that emits a range of sequential integers.
func Range(start, count int) ObservableInt {
	end := start + count
	observable := func(observe IntObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		i := start
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Subscribed() {
				if i < end {
					observe(i, nil, false)
					if subscriber.Subscribed() {
						i++
						self()
					}
				} else {
					observe(zeroInt, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template Observable Repeat

// Repeat creates an Observable that emits a sequence of items repeatedly.
func (o Observable) Repeat(count int) Observable {
	if count == 0 {
		return Empty()
	}
	observable := func(observe ObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		var repeated int
		var observer ObserveFunc
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

//jig:template Repeat<Foo>
//jig:needs Observable<Foo>

// RepeatFoo creates an ObservableFoo that emits a particular item or sequence
// of items repeatedly.
func RepeatFoo(value foo, count int) ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		i := 0
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Subscribed() {
				if i < count {
					observe(value, nil, false)
					if subscriber.Subscribed() {
						i++
						self()
					}
				} else {
					observe(zeroFoo, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template Start<Foo>
//jig:needs Observable<Foo>

// StartFoo creates an ObservableFoo that emits the return value of a function.
// It is designed to be used with a function that returns a (foo, error) tuple.
// If the error is non-nil the returned ObservableFoo will be an Observable that
// emits and error, otherwise it will be a single-value ObservableFoo of the value.
func StartFoo(f func() (foo, error)) ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		done := false
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Subscribed() {
				if !done {
					if next, err := f(); err == nil {
						observe(next, nil, false)
						if subscriber.Subscribed() {
							done = true
							self()
						}
					} else {
						observe(zeroFoo, err, true)
					}
				} else {
					observe(zeroFoo, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template Throw<Foo>
//jig:needs Observable<Foo>

// ThrowFoo creates an Observable that emits no items and terminates with an
// error.
func ThrowFoo(err error) ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Subscribed() {
				observe(zeroFoo, err, true)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}
