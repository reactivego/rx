package rx

import "sync/atomic"

//jig:template Observable All
//jig:needs ObservableBool

// All determines whether all items emitted by an Observable meet some
// criteria.
//
// Pass a predicate function to the All operator that accepts an item emitted
// by the source Observable and returns a boolean value based on an
// evaluation of that item. All returns an ObservableBool that emits a single
// boolean value: true if and only if the source Observable terminates
// normally and every item emitted by the source Observable evaluated as
// true according to this predicate; false if any item emitted by the source
// Observable evaluates as false according to this predicate.
func (o Observable) All(predicate func(next interface{}) bool) ObservableBool {
	observable := func(observe BoolObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			switch {
			case !done:
				if !predicate(next) {
					observe(false, nil, false)
					observe(false, nil, true)
				}
			case err!=nil:
				observe(false, err, true)
			default:
				observe(true, nil, false)
				observe(false, nil, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> All

// All determines whether all items emitted by an ObservableFoo meet some
// criteria.
//
// Pass a predicate function to the All operator that accepts an item emitted
// by the source ObservableFoo and returns a boolean value based on an
// evaluation of that item. All returns an ObservableBool that emits a single
// boolean value: true if and only if the source ObservableFoo terminates
// normally and every item emitted by the source ObservableFoo evaluated as
// true according to this predicate; false if any item emitted by the source
// ObservableFoo evaluates as false according to this predicate.
func (o ObservableFoo) All(predicate func(next foo) bool) ObservableBool {
	condition := func(next interface{}) bool {
		return predicate(next.(foo))
	}
	return o.AsObservable().All(condition)
}

//jig:template Observable TakeWhile

// TakeWhile mirrors items emitted by an Observable until a specified condition becomes false.
//
// The TakeWhile mirrors the source Observable until such time as some condition you specify
// becomes false, at which point TakeWhile stops mirroring the source Observable and terminates
// its own Observable.
func (o Observable) TakeWhile(condition func(next interface{}) bool) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if done || condition(next) {
				observe(next, err, done)
			} else {
				observe(nil, nil, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> TakeWhile
//jig:needs Observable TakeWhile

// TakeWhile mirrors items emitted by an Observable until a specified condition becomes false.
//
// The TakeWhile mirrors the source Observable until such time as some condition you specify
// becomes false, at which point TakeWhile stops mirroring the source Observable and terminates
// its own Observable.
func (o ObservableFoo) TakeWhile(condition func(next foo) bool) ObservableFoo {
	predecate := func (next interface{}) bool {
		return condition(next.(foo))
	}
	return o.AsObservable().TakeWhile(predecate).AsObservableFoo()
}

//jig:template Observable TakeUntil
//jig:needs nil

// TakeUntil emits items emitted by an Observable until another Observable emits an item.
func (o Observable) TakeUntil(other Observable) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var watcherNext int32
		watcherSubscriber := subscriber.Add()
		watcher := func (next interface{}, err error, done bool) {
			if !done {
				atomic.StoreInt32(&watcherNext, 1)
			}
			watcherSubscriber.Unsubscribe()
		}
		other(watcher, subscribeOn, watcherSubscriber)
		
		observer := func(next interface{}, err error, done bool) {
			if done || atomic.LoadInt32(&watcherNext) != 1 {
				observe(next, err, done)
			} else {
				observe(nil, nil, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> TakeUntil
//jig:needs Observable TakeUntil

// TakeUntil emits items emitted by an ObservableFoo until another Observable emits an item.
func (o ObservableFoo) TakeUntil(other Observable) ObservableFoo {
	return o.AsObservable().TakeUntil(other).AsObservableFoo()
}

//jig:template Observable<Foo> TakeUntil<Bar>
//jig:needs Observable<Foo> TakeUntil

// TakeUntil emits items emitted by an ObservableFoo until another ObservableBar emits an item.
func (o ObservableFoo) TakeUntilBar(other ObservableBar) ObservableFoo {
	return o.TakeUntil(other.AsObservable())
}
