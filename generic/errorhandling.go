package rx

//jig:template Observable<Foo> Catch

// Catch recovers from an error notification by continuing the sequence without
// emitting the error but by switching to the catch ObservableFoo to provide
// items.
func (o ObservableFoo) Catch(catch ObservableFoo) ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			if err != nil {
				catch(observe, subscribeOn, subscriber)
			} else {
				observe(next, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> CatchError

// CatchError catches errors on the ObservableFoo to be handled by returning a
// new ObservableFoo or throwing an error. It is passed a selector function 
// that takes as arguments err, which is the error, and caught, which is the
// source observable, in case you'd like to "retry" that observable by
// returning it again. Whatever observable is returned by the selector will be
// used to continue the observable chain.
func (o ObservableFoo) CatchError(selector func(err error, caught ObservableFoo) ObservableFoo) ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			if err != nil {
				selector(err, o)(observe, subscribeOn, subscriber)
			} else {
				observe(next, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Retry

// Retry if a source ObservableFoo sends an error notification, resubscribe to
// it in the hopes that it will complete without error.
func (o ObservableFoo) Retry() ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var observer FooObserver
		observer = func(next foo, err error, done bool) {
			if err != nil {
				o(observer, subscribeOn, subscriber)
			} else {
				observe(next, nil, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}
