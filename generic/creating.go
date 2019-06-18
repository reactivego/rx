package rx

import "time"

//jig:template <Foo>Observer

// Next is called by an ObservableFoo to emit the next foo value to the
// observer.
func (f FooObserveFunc) Next(next foo) {
	f(next, nil, false)
}

// Error is called by an ObservableFoo to report an error to the observer.
func (f FooObserveFunc) Error(err error) {
	f(zeroFoo, err, true)
}

// Complete is called by an ObservableFoo to signal that no more data is
// forthcoming to the observer.
func (f FooObserveFunc) Complete() {
	f(zeroFoo, nil, true)
}

// FooObserver is the interface used with CreateFoo when implementing a custom
// observable.
type FooObserver interface {
	// Next emits the next foo value.
	Next(foo)
	// Error signals an error condition.
	Error(error)
	// Complete signals that no more data is to be expected.
	Complete()
	// Closed returns true when the subscription has been canceled.
	Closed() bool
}

//jig:template Create<Foo>
//jig:needs Observable<Foo>, <Foo>Observer

// CreateFoo creates an Observable from scratch by calling observer methods
// programmatically.
func CreateFoo(f func(FooObserver)) ObservableFoo {
	observable := func(observe FooObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		scheduler.Schedule(func() {
			if subscriber.Closed() {
				return
			}
			observer := func(next foo, err error, done bool) {
				if !subscriber.Closed() {
					observe(next, err, done)
				}
			}
			type ObserverSubscriber struct {
				FooObserveFunc
				Subscriber
			}
			f(&ObserverSubscriber{observer, subscriber})
		})
	}
	return observable
}

//jig:template Defer<Foo>
//jig:needs Create<Foo>

// DeferFoo does not create the ObservableFoo until the observer subscribes,
// and creates a fresh ObservableFoo for each observer.
func DeferFoo(factory func() ObservableFoo) ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		factory()(observe, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Empty<Foo>
//jig:needs Create<Foo>

// EmptyFoo creates an Observable that emits no items but terminates normally.
func EmptyFoo() ObservableFoo {
	return CreateFoo(func(observer FooObserver) {
		observer.Complete()
	})
}

//jig:template Error<Foo>
//jig:needs Observable<Foo>

// ErrorFoo creates an Observable that emits no items and terminates with an
// error.
func ErrorFoo(err error) ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		subscribeOn.Schedule(func() {
			if !subscriber.Canceled() {
				observe(zeroFoo, err, true)
			}
		})
	}
	return observable
}

//jig:template FromChan<Foo>
//jig:needs Create<Foo>

// FromChanFoo creates an ObservableFoo from a Go channel of foo values.
// It's not possible for the code feeding into the channel to send an error.
// The feeding code can send zero or more foo items and then closing the
// channel will be seen as completion.
func FromChanFoo(ch <-chan foo) ObservableFoo {
	return CreateFoo(func(observer FooObserver) {
		for next := range ch {
			if observer.Closed() {
				return
			}
			observer.Next(next)
		}
		observer.Complete()
	})
}

//jig:template FromChan
//jig:needs Create

// FromChan creates an Observable from a Go channel of interface{}
// values. This allows the code feeding into the channel to send either an error
// or the next value. The feeding code can send zero or more items and then
// closing the channel will be seen as completion. When the feeding code sends
// an error into the channel, it should close the channel immediately to
// indicate termination with error.
func FromChan(ch <-chan interface{}) Observable {
	return Create(func(observer Observer) {
		for item := range ch {
			if observer.Closed() {
				return
			}
			if err, ok := item.(error); ok {
				observer.Error(err)
				return
			}
			observer.Next(item)
		}
		if observer.Closed() {
			return
		}
		observer.Complete()
	})
}

//jig:template From<Foo>
//jig:needs FromSlice<Foo>

// FromFoo creates an ObservableFoo from multiple foo values passed in.
func FromFoo(slice ...foo) ObservableFoo {
	return FromSliceFoo(slice)
}

//jig:template From<Foo>s
//jig:needs FromSlice<Foo>

// FromFoos creates an ObservableFoo from multiple foo values passed in.
func FromFoos(slice ...foo) ObservableFoo {
	return FromSliceFoo(slice)
}

//jig:template FromSlice<Foo>
//jig:needs Observable<Foo>

// FromSliceFoo creates an ObservableFoo from a slice of foo values passed in.
func FromSliceFoo(slice []foo) ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		subscribeOn.ScheduleRecursive(func(self func()) {
			if !subscriber.Canceled() {
				if i < len(slice) {
					observe(slice[i], nil, false)
					if !subscriber.Canceled() {
						i++
						self()
					}
				} else {
					observe(zeroFoo, nil, true)
				}
			}
		})
	}
	return observable
}

//jig:template Interval
//jig:needs CreateInt

// Interval creates an ObservableInt that emits a sequence of integers spaced
// by a particular time interval.
func Interval(interval time.Duration) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		for i := 0; ; i++ {
			time.Sleep(interval)
			if observer.Closed() {
				return
			}
			observer.Next(i)
		}
	})
}

//jig:template Just<Foo>
//jig:needs Create<Foo>

// JustFoo creates an ObservableFoo that emits a particular item.
func JustFoo(element foo) ObservableFoo {
	return CreateFoo(func(observer FooObserver) {
		observer.Next(element)
		observer.Complete()
	})
}

//jig:template Never<Foo>
//jig:needs Create<Foo>

// NeverFoo creates an ObservableFoo that emits no items and does't terminate.
func NeverFoo() ObservableFoo {
	return CreateFoo(func(observer FooObserver) {
	})
}

//jig:template Range
//jig:needs CreateInt

// Range creates an ObservableInt that emits a range of sequential integers.
func Range(start, count int) ObservableInt {
	end := start + count
	observable := func(observe IntObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		i := start
		scheduler.ScheduleRecursive(func(self func()) {
			if !subscriber.Closed() {
				if i < end {
					observe(i, nil, false)
					if !subscriber.Closed() {
						i++
						self()
					}
				} else {
					observe(zeroInt, nil, true)
				}
			}
		})
	}
	return observable
}

//jig:template Observable Repeat

// Repeat creates an Observable that emits a sequence of items repeatedly.
func (o Observable) Repeat(count int) Observable {
	if count == 0 {
		return Empty()
	}
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var repeated int
		var observer ObserveFunc
		observer = func(next interface{}, err error, done bool) {
			if !done || err != nil {
				observe(next, err, done)
			} else {
				repeated++
				if repeated < count {
					o(observer, subscribeOn, subscriber)
				} else {
					observe(nil, nil, true)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
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
//jig:needs Create<Foo>

// RepeatFoo creates an ObservableFoo that emits a particular item or sequence
// of items repeatedly.
func RepeatFoo(value foo, count int) ObservableFoo {
	return CreateFoo(func(observer FooObserver) {
		for i := 0; i < count; i++ {
			if observer.Closed() {
				return
			}
			observer.Next(value)
		}
		observer.Complete()
	})
}

//jig:template Start<Foo>
//jig:needs Create<Foo>

// StartFoo creates an ObservableFoo that emits the return value of a function.
// It is designed to be used with a function that returns a (foo, error) tuple.
// If the error is non-nil the returned ObservableFoo will be that error,
// otherwise it will be a single-value stream of foo.
func StartFoo(f func() (foo, error)) ObservableFoo {
	return CreateFoo(func(observer FooObserver) {
		if next, err := f(); err == nil {
			observer.Next(next)
			observer.Complete()
		} else {
			observer.Error(err)
		}
	})
}

//jig:template Throw<Foo>
//jig:needs Create<Foo>

// ThrowFoo creates an Observable that emits no items and terminates with an
// error.
func ThrowFoo(err error) ObservableFoo {
	return CreateFoo(func(observer FooObserver) {
		observer.Error(err)
	})
}
