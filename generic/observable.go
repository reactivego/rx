package rx

//jig:template <Foo>ObserveFunc

// FooObserveFunc is the observer, a function that gets called whenever the
// observable has something to report. The next argument is the item value that
// is only valid when the done argument is false. When done is true and the err
// argument is not nil, then the observable has terminated with an error.
// When done is true and the err argument is nil, then the observable has
// completed normally.
type FooObserveFunc func(next foo, err error, done bool)

var zeroFoo foo

// Next is called by an ObservableFoo to emit the next foo value to the
// observer.
func (f FooObserveFunc) Next(next foo) {
	f(next, nil, false)
}

// Error is called by an ObservableFoo to report an error to the observer.
func (f FooObserveFunc) Error(err error) {
	f(zeroFoo, err, true)
}

// Complete is called by an ObservableFoo to signal that no more data is
// forthcoming to the observer.
func (f FooObserveFunc) Complete() {
	f(zeroFoo, nil, true)
}

//jig:template Observable<Foo>
//jig:needs <Foo>ObserveFunc

// ObservableFoo is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableFoo func(FooObserveFunc, Scheduler, Subscriber)
