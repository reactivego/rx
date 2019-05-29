package rx

import "sync/atomic"

//jig:template Observable TakeWhile

// TakeWhile mirrors items emitted by an Observable until a specified condition becomes false.
//
// The TakeWhile mirrors the source Observable until such time as some condition you specify
// becomes false, at which point TakeWhile stops mirroring the source Observable and terminates
// its own Observable.
func (o Observable) TakeWhile(condition func(next interface{}) bool) Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if done || condition(next) {
				observe(next, err, done)
			} else {
				observe(zero, nil, true)
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

// TakeUntil emits items emitted by an Observable until another Observable emits an item.
func (o Observable) TakeUntil(other Observable) Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var watcherNext int32
		watcherSubscriber := subscriber.Add(func(){})
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
				observe(zero, nil, true)
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
