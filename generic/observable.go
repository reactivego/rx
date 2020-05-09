package rx

//jig:template zero<Foo>

var zeroFoo foo

//jig:template <Foo>Observer

// FooObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type FooObserver func(next foo, err error, done bool)

//jig:template Observable<Foo>
//jig:needs Subscriber, Scheduler, <Foo>Observer

// ObservableFoo is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableFoo func(FooObserver, Scheduler, Subscriber)
