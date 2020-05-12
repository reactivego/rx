package rx

import (
	"sync"
	"sync/atomic"
)

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
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
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
					var zeroFoo foo
					observe(zeroFoo, err, true)
				default:
					if observers.len--; observers.len == 0 {
						var zeroFoo foo
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
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
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
					var zeroFoo foo
					observe(zeroFoo, err, true)
				default:
					if atomic.AddInt32(&observers.len, -1) == 0 {
						var zeroFoo foo
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
				var zeroFoo foo
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
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
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
					var zeroFoo foo
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
