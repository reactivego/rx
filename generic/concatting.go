package rx

import "sync"

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

//jig:template ObservableObservable<Foo> ConcatAll

// ConcatAll flattens a higher order observable by concattenating the observables it emits.
func (o ObservableObservableFoo) ConcatAll() ObservableFoo {
	var zeroFoo foo
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			mutex       sync.Mutex
			observables []ObservableFoo
			observer    FooObserver
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
		sourceSubscriber := subscriber.Add()
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
