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
