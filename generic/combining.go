package rx

import (
	"sync"
	"sync/atomic"
)

//jig:template ObservableObservable<Foo> CombineAll

type FooSlice []foo

func (o ObservableObservableFoo) CombineAll() ObservableFooSlice {
	observable := func(observe FooSliceObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observables := []ObservableFoo(nil)
		var observers struct {
			sync.Mutex
			values []foo
			initialized int
			active int
		}
		makeObserver := func(index int) FooObserveFunc {
			observer := func(next foo, err error, done bool) {
				observers.Lock()
				defer observers.Unlock()
				if observers.active > 0 {
					switch {
					case !done:
						if observers.values[index] == zeroFoo {
							observers.initialized++
						}
						observers.values[index] = next
						if observers.initialized == len(observers.values) {
							observe(observers.values, nil, false)
						}
					case err != nil:
						observers.active = 0
						observe(zeroFooSlice, err, true)
					default:
						if observers.active--; observers.active == 0 {
							observe(zeroFooSlice, nil, true)
						}
					}
				}
			}
			return observer
		}

		observer := func(next ObservableFoo, err error, done bool) {
			switch {
			case !done:
				observables = append(observables, next)
			case err != nil:
				observe(zeroFooSlice, err, true)
			default:
				subscribeOn.Schedule(func() {
					if !subscriber.Canceled() {
						numObservables := len(observables)
						observers.values = make([]foo, numObservables)
						observers.active = numObservables
						for i, v := range observables {
							if subscriber.Canceled() {
								return
							}
							v(makeObserver(i), subscribeOn, subscriber)
						}
					}
				})
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Concat<Foo>
//jig:needs Observable<Foo> Concat

// ConcatFoo emits the emissions from two or more ObservableFoos without interleaving them.
func ConcatFoo(observables ...ObservableFoo) ObservableFoo {
	if len(observables) == 0 {
		return EmptyFoo()
	}
	return observables[0].Concat(observables[1:]...)
}

//jig:template Observable<Foo> Concat

// Concat emits the emissions from two or more ObservableFoos without interleaving them.
func (o ObservableFoo) Concat(other ...ObservableFoo) ObservableFoo {
	if len(other) == 0 {
		return o
	}
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			observables = append([]ObservableFoo{}, other...)
			observer    FooObserveFunc
		)
		observer = func(next foo, err error, done bool) {
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if len(observables) == 0 {
					observe(zeroFoo, nil, true)
				} else {
					o := observables[0]
					observables = observables[1:]
					o(observer, subscribeOn, subscriber)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template ObservableObservable<Foo> ConcatAll

// ConcatAll flattens a higher order observable by concattenating the observables it emits.
func (o ObservableObservableFoo) ConcatAll() ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			mutex       sync.Mutex
			observables []ObservableFoo
			observer    FooObserveFunc
		)
		observer = func(next foo, err error, done bool) {
			mutex.Lock()
			defer mutex.Unlock()
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if len(observables) == 0 {
					observe(zeroFoo, nil, true)
				} else {
					o := observables[0]
					observables = observables[1:]
					o(observer, subscribeOn, subscriber)
				}
			}
		}
		sourceSubscriber := subscriber.AddChild()
		concatenator := func(next ObservableFoo, err error, done bool) {
			if !done {
				mutex.Lock()
				defer mutex.Unlock()
				observables = append(observables, next)
			} else {
				observer(zeroFoo, err, done)
				sourceSubscriber.Unsubscribe()
			}
		}
		o(concatenator, subscribeOn, sourceSubscriber)
	}
	return observable
}

//jig:template Merge<Foo>
//jig:needs Observable<Foo> Merge

// MergeFoo combines multiple Observables into one by merging their emissions.
// An error from any of the observables will terminate the merged observables.
func MergeFoo(observables ...ObservableFoo) ObservableFoo {
	if len(observables) == 0 {
		return EmptyFoo()
	}
	return observables[0].Merge(observables[1:]...)
}

//jig:template Observable<Foo> Merge

// Merge combines multiple Observables into one by merging their emissions.
// An error from any of the observables will terminate the merged observables.
func (o ObservableFoo) Merge(other ...ObservableFoo) ObservableFoo {
	if len(other) == 0 {
		return o
	}
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var observers struct {
			sync.Mutex
			done bool
			len  int
		}
		observer := func(next foo, err error, done bool) {
			observers.Lock()
			defer observers.Unlock()
			if !observers.done {
				switch {
				case !done:
					observe(next, nil, false)
				case err != nil:
					observers.done = true
					observe(zeroFoo, err, true)
				default:
					if observers.len--; observers.len == 0 {
						observe(zeroFoo, nil, true)
					}
				}
			}
		}
		subscribeOn.Schedule(func() {
			if !subscriber.Canceled() {
				observers.len = 1 + len(other)
				o(observer, subscribeOn, subscriber)
				for _, o := range other {
					if subscriber.Canceled() {
						return
					}
					o(observer, subscribeOn, subscriber)
				}
			}
		})
	}
	return observable
}

//jig:template ObservableObservable<Foo> MergeAll

// MergeAll flattens a higher order observable by merging the observables it emits.
func (o ObservableObservableFoo) MergeAll() ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var observers struct {
			sync.Mutex
			done bool
			len  int32
		}
		observer := func(next foo, err error, done bool) {
			observers.Lock()
			defer observers.Unlock()
			if !observers.done {
				switch {
				case !done:
					observe(next, nil, false)
				case err != nil:
					observers.done = true
					observe(zeroFoo, err, true)
				default:
					if atomic.AddInt32(&observers.len, -1) == 0 {
						observe(zeroFoo, nil, true)
					}
				}
			}
		}
		merger := func(next ObservableFoo, err error, done bool) {
			if !done {
				atomic.AddInt32(&observers.len, 1)
				next(observer, subscribeOn, subscriber)
			} else {
				observer(zeroFoo, err, true)
			}
		}
		subscribeOn.Schedule(func() {
			if !subscriber.Canceled() {
				observers.len = 1
				o(merger, subscribeOn, subscriber)
			}
		})
	}
	return observable
}

//jig:template MergeDelayError<Foo>
//jig:needs Observable<Foo> MergeDelayError

// MergeDelayErrorFoo combines multiple Observables into one by merging their emissions.
// Any error will be deferred until all observables terminate.
func MergeDelayErrorFoo(observables ...ObservableFoo) ObservableFoo {
	if len(observables) == 0 {
		return EmptyFoo()
	}
	return observables[0].MergeDelayError(observables[1:]...)
}

//jig:template Observable<Foo> MergeDelayError

// MergeDelayError combines multiple Observables into one by merging their emissions.
// Any error will be deferred until all observables terminate.
func (o ObservableFoo) MergeDelayError(other ...ObservableFoo) ObservableFoo {
	if len(other) == 0 {
		return o
	}
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var observers struct {
			sync.Mutex
			len int
			err error
		}
		observer := func(next foo, err error, done bool) {
			observers.Lock()
			defer observers.Unlock()
			if !done {
				observe(next, nil, false)
			} else {
				if err != nil {
					observers.err = err
				}
				if observers.len--; observers.len == 0 {
					observe(zeroFoo, observers.err, true)
				}
			}
		}
		subscribeOn.Schedule(func() {
			if !subscriber.Canceled() {
				observers.len = 1 + len(other)
				o(observer, subscribeOn, subscriber)
				for _, o := range other {
					if subscriber.Canceled() {
						return
					}
					o(observer, subscribeOn, subscriber)
				}
			}
		})
	}
	return observable
}

//jig:template ObservableObservable<Foo> SwitchAll
//jig:needs <Foo>Link

// SwitchAll converts an Observable that emits Observables into a single Observable
// that emits the items emitted by the most-recently-emitted of those Observables.
func (o ObservableObservableFoo) SwitchAll() ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(link *FooLink, next foo, err error, done bool) {
			if !done || err != nil {
				observe(next, err, done)
			} else {
				link.subscriber.Unsubscribe() // We filter complete. Therefore, we need to perform Unsubscribe.
			}
		}
		currentLink := NewInitialFooLink()
		var switcherMutex sync.Mutex
		switcherSubscriber := subscriber.AddChild()
		switcher := func(next ObservableFoo, err error, done bool) {
			switch {
			case !done:
				previousLink := currentLink
				func() {
					switcherMutex.Lock()
					defer switcherMutex.Unlock()
					currentLink = NewFooLink(observer, subscriber)
				}()
				previousLink.Cancel(func() {
					switcherMutex.Lock()
					defer switcherMutex.Unlock()
					currentLink.SubscribeTo(next, subscribeOn)
				})
			case err != nil:
				currentLink.Cancel(func() {
					observe(zeroFoo, err, true)
				})
				switcherSubscriber.Unsubscribe()
			default:
				currentLink.OnComplete(func() {
					observe(zeroFoo, nil, true)
				})
				switcherSubscriber.Unsubscribe()
			}
		}
		o(switcher, subscribeOn, switcherSubscriber)
	}
	return observable
}
