package rx

import (
	"sync"
	"sync/atomic"
)

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
		sourceSubscriber := subscriber.Add(func() { /*unsubscribed*/ })
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
		var (
			mutex sync.Mutex
			count = 1 + len(other)
		)
		observer := func(next foo, err error, done bool) {
			mutex.Lock()
			defer mutex.Unlock()
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if count--; count == 0 {
					observe(zeroFoo, nil, true)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
		for _, o := range other {
			o(observer, subscribeOn, subscriber)
		}
	}
	return observable
}

//jig:template ObservableObservable<Foo> MergeAll

// MergeAll flattens a higher order observable by merging the observables it emits.
func (o ObservableObservableFoo) MergeAll() ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			mutex sync.Mutex
			count int32 = 1
		)
		observer := func(next foo, err error, done bool) {
			mutex.Lock()
			defer mutex.Unlock()
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if atomic.AddInt32(&count, -1) == 0 {
					observe(zeroFoo, nil, true)
				}
			}
		}
		merger := func(next ObservableFoo, err error, done bool) {
			if !done {
				atomic.AddInt32(&count, 1)
				next(observer, subscribeOn, subscriber)
			} else {
				observer(zeroFoo, err, true)
			}
		}
		o(merger, subscribeOn, subscriber)
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
		var (
			mutex      sync.Mutex
			count      = 1 + len(other)
			delayedErr error
		)
		observer := func(next foo, err error, done bool) {
			mutex.Lock()
			defer mutex.Unlock()
			if !done {
				observe(next, nil, false)
			} else {
				if err != nil {
					delayedErr = err
				}
				if count--; count == 0 {
					observe(zeroFoo, delayedErr, true)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
		for _, o := range other {
			o(observer, subscribeOn, subscriber)
		}
	}
	return observable
}
