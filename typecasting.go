package rx

import "errors"

//jig:template Observable<Foo> AsAny
//jig:needs Observable

// AsAny turns a typed ObservableFoo into an Observable of interface{}.
func (o ObservableFoo) AsAny() Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			observe(next, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template ErrTypecastTo<Foo>

// ErrTypecastToFoo is delivered to an observer if the generic value cannot be
// typecast to foo.
var ErrTypecastToFoo = errors.New("typecast to foo failed")

//jig:template Observable As<Foo>
//jig:needs <Foo>ObserveFunc, ErrTypecastTo<Foo>
//jig:required-vars Foo

// AsFoo turns an Observable of interface{} into an ObservableFoo. If during
// observing a typecast fails, the error ErrTypecastToFoo will be emitted.
func (o Observable) AsFoo() ObservableFoo {
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

//jig:template Observable<Foo> As<Bar>
//jig:needs Observable<Foo>
//jig:required-vars Foo, Bar

// AsAny turns a typed ObservableFoo into an Observable of interface{}.
func (o ObservableFoo) AsBar() ObservableBar {
	observable := func(observe BarObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			observe(bar(next), err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable Only<Foo>
//jig:needs <Foo>ObserveFunc
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
