package rx

import "sync"

//jig:template Observable<Foo> WithLatestFrom
//jig:needs ObservableObservable<Foo> WithLatestFromAll

// WithLatestFrom will subscribe to all Observables and wait for all of them to emit before emitting
// the first slice. The source observable determines the rate at which the values are emitted. The idea
// is that observables that are faster than the source, don't determine the rate at which the resulting
// observable emits. The observables that are combined with the source will be allowed to continue
// emitting but only will have their last emitted value emitted whenever the source emits.
func (o ObservableFoo) WithLatestFrom(other ...ObservableFoo) ObservableFooSlice {
	return FromObservableFoo(append([]ObservableFoo{o}, other...)...).WithLatestFromAll()
}


//jig:template ObservableObservable<Foo> WithLatestFromAll
//jig:needs <Foo>Slice

// WithLatestFromAll flattens a higher order observable (e.g. ObservableObservableFoo) by
// subscribing to all emitted observables (ie. ObservableFoo entries) until the source completes. It
// will then wait for all of the subscribed ObservableFoos to emit before emitting the first slice.
// Whenever the first emitted observable emits, a new slice will be emitted containing all the
// latest value.
func (o ObservableObservableFoo) WithLatestFromAll() ObservableFooSlice {
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
						if index == 0 && observers.initialized == len(observers.values){
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