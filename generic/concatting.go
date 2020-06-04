package rx

import "sync"

//jig:template Concat<Foo>
//jig:needs Observable<Foo> Concat

// ConcatFoo emits the emissions from two or more ObservableFoos without interleaving them.
func ConcatFoo(observables ...ObservableFoo) ObservableFoo {
	if len(observables) == 0 {
		return EmptyFoo()
	}
	return observables[0].ConcatWith(observables[1:]...)
}

//jig:template Observable<Foo> ConcatWith

// ConcatWith emits the emissions from two or more ObservableFoos without interleaving them.
func (o ObservableFoo) ConcatWith(other ...ObservableFoo) ObservableFoo {
	var zeroFoo foo
	if len(other) == 0 {
		return o
	}
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			observables = append([]ObservableFoo{}, other...)
			observer    FooObserver
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

//jig:template Observable<Foo> ConcatMap<Bar>

// ConcatMapBar transforms the items emitted by an ObservableFoo by applying a
// function to each item and returning an ObservableBar. The stream of
// ObservableBar items is then flattened by concattenating the emissions from
// the observables without interleaving.
func (o ObservableFoo) ConcatMapBar(project func(foo) ObservableBar) ObservableBar {
	return o.MapObservableBar(project).ConcatAll()
}

//jig:template Observable<Foo> ConcatMapTo<Bar>

// ConcatMapToBar maps every entry emitted by the ObservableFoo into a single
// ObservableBar. The stream of ObservableBar items is then flattened by
// concattenating the emissions from the observables without interleaving.
func (o ObservableFoo) ConcatMapToBar(inner ObservableBar) ObservableBar {
	project := func(foo) ObservableBar { return inner }
	return o.MapObservableBar(project).ConcatAll()
}

//jig:template ObservableObservable<Foo> ConcatAll

// ConcatAll flattens a higher order observable by concattenating the observables it emits.
func (o ObservableObservableFoo) ConcatAll() ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var concat struct {
			sync.Mutex
			observables []ObservableFoo
			observer    FooObserver
			subscriber  Subscriber
		}

		var source struct {
			observer   ObservableFooObserver
			subscriber Subscriber
		}

		concat.observer = func(next foo, err error, done bool) {
			concat.Lock()
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if len(concat.observables) == 0 {
					if !source.subscriber.Subscribed() {
						var zero foo
						observe(zero, nil, true)
					}
					concat.observables = nil
				} else {
					observable := concat.observables[0]
					concat.observables = concat.observables[1:]
					observable(concat.observer, subscribeOn, subscriber)
				}
			}
			concat.Unlock()
		}

		source.observer = func(next ObservableFoo, err error, done bool) {
			if !done {
				concat.Lock()
				initial := concat.observables == nil
				concat.observables = append(concat.observables, next)
				concat.Unlock()
				if initial {
					var zero foo
					concat.observer(zero, nil, true)
				}
			} else {
				concat.Lock()
				initial := concat.observables == nil
				source.subscriber.Done(err)
				concat.Unlock()
				if initial || err != nil {
					var zero foo
					concat.observer(zero, err, true)
				}
			}
		}
		source.subscriber = subscriber.Add()
		o(source.observer, subscribeOn, source.subscriber)
	}
	return observable
}
