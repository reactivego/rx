package rx

//jig:template <Foo>ObserveFunc

// FooObserveFunc is the observer, a function that gets called whenever the
// observable has something to report. The next argument is the item value that
// is only valid when the done argument is false. When done is true and the err
// argument is not nil, then the observable has terminated with an error.
// When done is true and the err argument is nil, then the observable has
// completed normally.
type FooObserveFunc func(next foo, err error, done bool)

//jig:template zero<Foo>

var zeroFoo foo

//jig:template Observable<Foo>
//jig:needs Scheduler, Subscriber, <Foo>ObserveFunc, zero<Foo>

// ObservableFoo is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableFoo func(FooObserveFunc, Scheduler, Subscriber)
