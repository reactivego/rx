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
	return observables[0].MergeWith(observables[1:]...)
}

//jig:template Observable<Foo> MergeWith

// MergeWith combines multiple Observables into one by merging their emissions.
// An error from any of the observables will terminate the merged observables.
func (o ObservableFoo) MergeWith(other ...ObservableFoo) ObservableFoo {
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
					var zero foo
					observe(zero, err, true)
				default:
					if observers.len--; observers.len == 0 {
						var zero foo
						observe(zero, nil, true)
					}
				}
			}
		}
		subscribeOn.Schedule(func() {
			if subscriber.Subscribed() {
				observers.len = 1 + len(other)
				o(observer, subscribeOn, subscriber)
				for _, o := range other {
					if !subscriber.Subscribed() {
						return
					}
					o(observer, subscribeOn, subscriber)
				}
			}
		})
	}
	return observable
}

//jig:template Observable<Foo> MergeMap<Bar>

// MergeMapBar transforms the items emitted by an ObservableFoo by applying a
// function to each item an returning an ObservableBar. The stream of ObservableBar
// items is then merged into a single stream of Bar items using the MergeAll operator.
func (o ObservableFoo) MergeMapBar(project func(foo) ObservableBar) ObservableBar {
	return o.MapObservableBar(project).MergeAll()
}

//jig:template Observable<Foo> MergeMapTo<Bar>

// MergeMapToBar maps every entry emitted by the ObservableFoo into a single
// ObservableBar. The stream of ObservableBar items is then merged into a
// single stream of Bar items using the MergeAll operator.
func (o ObservableFoo) MergeMapToBar(inner ObservableBar) ObservableBar {
	project := func(foo) ObservableBar { return inner }
	return o.MapObservableBar(project).MergeAll()
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
					var zero foo
					observe(zero, err, true)
				default:
					if atomic.AddInt32(&observers.len, -1) == 0 {
						var zero foo
						observe(zero, nil, true)
					}
				}
			}
		}
		merger := func(next ObservableFoo, err error, done bool) {
			if !done {
				atomic.AddInt32(&observers.len, 1)
				next(observer, subscribeOn, subscriber)
			} else {
				var zero foo
				observer(zero, err, true)
			}
		}
		subscribeOn.Schedule(func() {
			if subscriber.Subscribed() {
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
	return observables[0].MergeDelayErrorWith(observables[1:]...)
}

//jig:template Observable<Foo> MergeDelayErrorWith

// MergeDelayError combines multiple Observables into one by merging their emissions.
// Any error will be deferred until all observables terminate.
func (o ObservableFoo) MergeDelayErrorWith(other ...ObservableFoo) ObservableFoo {
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
					var zero foo
					observe(zero, observers.err, true)
				}
			}
		}
		subscribeOn.Schedule(func() {
			if subscriber.Subscribed() {
				observers.len = 1 + len(other)
				o(observer, subscribeOn, subscriber)
				for _, o := range other {
					if !subscriber.Subscribed() {
						return
					}
					o(observer, subscribeOn, subscriber)
				}
			}
		})
	}
	return observable
}
