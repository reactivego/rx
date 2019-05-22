package rx

//jig:template <Foo>ObserveFunc

// FooObserveFunc is essentially the observer, a function that gets called
// whenever the observable has something to report.
type FooObserveFunc func(foo, error, bool)

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
