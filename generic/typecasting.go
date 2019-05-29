package rx

//jig:template Observable<Foo> AsObservable<Bar>
//jig:needs Observable<Bar>
//jig:required-vars Foo

// AsObservableBar turns a typed ObservableFoo into an Observable of bar.
func (o ObservableFoo) AsObservableBar() ObservableBar {
	observable := func(observe BarObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			observe(bar(next), err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template ErrTypecastTo<Foo>
//jig:needs RxError

// ErrTypecastToFoo is delivered to an observer if the generic value cannot be
// typecast to foo.
const ErrTypecastToFoo = RxError("typecast to foo failed")

//jig:template Observable AsObservable<Foo>
//jig:needs Observable<Foo>, ErrTypecastTo<Foo>
//jig:required-vars Foo

// AsFoo turns an Observable of interface{} into an ObservableFoo. If during
// observing a typecast fails, the error ErrTypecastToFoo will be emitted.
func (o Observable) AsObservableFoo() ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if !done {
				if nextFoo, ok := next.(foo); ok {
					observe(nextFoo, err, done)
				} else {
					observe(zeroFoo, ErrTypecastToFoo, true)
				}
			} else {
				observe(zeroFoo, err, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable Only<Foo>
//jig:needs Observable<Foo>
//jig:required-vars Foo

// OnlyFoo filters the value stream of an Observable of interface{} and outputs only the
// foo typed values.
func (o Observable) OnlyFoo() ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if !done {
				if nextFoo, ok := next.(foo); ok {
					observe(nextFoo, err, done)
				}
			} else {
				observe(zeroFoo, err, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}
