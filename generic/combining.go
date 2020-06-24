package rx

import (
	"sync"
)

//jig:template CombineLatest<Foo>
//jig:needs ObservableObservable<Foo> CombineLatestAll

// CombineLatest will subscribe to all ObservableFoos. It will then wait for
// all of them to emit before emitting the first slice. Whenever any of the
// subscribed observables emits, a new slice will be emitted containing all
// the latest value.
func CombineLatestFoo(observables ...ObservableFoo) ObservableFooSlice {
	return FromObservableFoo(observables...).CombineLatestAll()
}

//jig:template Observable<Foo> CombineLatestWith
//jig:needs ObservableObservable<Foo> CombineLatestAll

// CombineLatestWith will subscribe to its ObservableFoo and all other
// ObservableFoos passed in. It will then wait for all of the ObservableBars
// to emit before emitting the first slice. Whenever any of the subscribed
// observables emits, a new slice will be emitted containing all the latest
// value.
func (o ObservableFoo) CombineLatestWith(other ...ObservableFoo) ObservableFooSlice {
	return FromObservableFoo(append([]ObservableFoo{o}, other...)...).CombineLatestAll()
}

//jig:template Observable<Foo> CombineLatestMap<Bar>
//jig:needs Observable<Foo> MapObservable<Bar>, ObservableObservable<Bar> CombineLatestAll

// CombinesLatestMap maps every entry emitted by the ObservableFoo into an
// ObservableBar, and then subscribe to it, until the source observable
// completes. It will then wait for all of the ObservableBars to emit before
// emitting the first slice. Whenever any of the subscribed observables emits,
// a new slice will be emitted containing all the latest value.
func (o ObservableFoo) CombineLatestMapBar(project func(foo) ObservableBar) ObservableBarSlice {
	return o.MapObservableBar(project).CombineLatestAll()
}

//jig:template Observable<Foo> CombineLatestMapTo<Bar>
//jig:needs Observable<Foo> MapObservable<Bar>, ObservableObservable<Bar> CombineLatestAll

// CombinesLatestMapTo maps every entry emitted by the ObservableFoo into a
// single ObservableBar, and then subscribe to it, until the source
// observable completes. It will then wait for all of the ObservableBars
// to emit before emitting the first slice. Whenever any of the subscribed
// observables emits, a new slice will be emitted containing all the latest
// value.
func (o ObservableFoo) CombineLatestMapToBar(inner ObservableBar) ObservableBarSlice {
	project := func(foo) ObservableBar { return inner }
	return o.MapObservableBar(project).CombineLatestAll()
}

//jig:template <Foo>Slice

type FooSlice = []foo

//jig:template ObservableObservable<Foo> CombineLatestAll
//jig:needs <Foo>Slice

// CombineLatestAll flattens a higher order observable
// (e.g. ObservableObservableFoo) by subscribing to
// all emitted observables (ie. ObservableFoo entries) until the source
// completes. It will then wait for all of the subscribed ObservableFoos
// to emit before emitting the first slice. Whenever any of the subscribed
// observables emits, a new slice will be emitted containing all the latest
// value.
func (o ObservableObservableFoo) CombineLatestAll() ObservableFooSlice {
	observable := func(observe FooSliceObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observables := []ObservableFoo(nil)
		var observers struct {
			sync.Mutex
			assigned    []bool
			values      []foo
			initialized int
			active      int
		}
		makeObserver := func(index int) FooObserver {
			observer := func(next foo, err error, done bool) {
				observers.Lock()
				defer observers.Unlock()
				if observers.active > 0 {
					switch {
					case !done:
						if !observers.assigned[index] {
							observers.assigned[index] = true
							observers.initialized++
						}
						observers.values[index] = next
						if observers.initialized == len(observers.values) {
							observe(observers.values, nil, false)
						}
					case err != nil:
						observers.active = 0
						var zero []foo
						observe(zero, err, true)
					default:
						if observers.active--; observers.active == 0 {
							var zero []foo
							observe(zero, nil, true)
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
				var zero []foo
				observe(zero, err, true)
			default:
				subscribeOn.Schedule(func() {
					if subscriber.Subscribed() {
						numObservables := len(observables)
						observers.assigned = make([]bool, numObservables)
						observers.values = make([]foo, numObservables)
						observers.active = numObservables
						for i, v := range observables {
							if !subscriber.Subscribed() {
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
