package rx

import (
	"time"

	"github.com/reactivego/multicast"
)


//jig:template Subject<Foo>
//jig:embeds Observable<Foo>, <Foo>Observer

// SubjectFoo is a combination of an observer and observable. Subjects are
// special because they are the only reactive constructs that support
// multicasting. The items sent to it through its observer side are
// multicasted to multiple clients subscribed to its observable side.
//
// A SubjectFoo embeds ObservableFoo and FooObserver. This exposes the
// methods and fields of both types on SubjectFoo. Use the ObservableFoo
// methods to subscribe to it. Use the FooObserver Next, Error and Complete
// methods to feed data to it.
//
// After a subject has been terminated by calling either Error or Complete,
// it goes into terminated state. All subsequent calls to its observer side
// will be silently ignored. All subsequent subscriptions to the observable
// side will be handled according to the specific behavior of the subject.
// There are different types of subjects, see the different NewXxxSubjectFoo
// functions for more info.
type SubjectFoo struct {
	ObservableFoo
	FooObserver
}

// Next is called by an ObservableFoo to emit the next foo value to the
// Observer.
func (f FooObserver) Next(next foo) {
	f(next, nil, false)
}

// Error is called by an ObservableFoo to report an error to the Observer.
func (f FooObserver) Error(err error) {
	var zeroFoo foo
	f(zeroFoo, err, true)
}

// Complete is called by an ObservableFoo to signal that no more data is
// forthcoming to the Observer.
func (f FooObserver) Complete() {
	var zeroFoo foo
	f(zeroFoo, nil, true)
}

//jig:template NewSubject<Foo>
//jig:needs Subject<Foo>

// NewSubjectFoo creates a new Subject. After the subject is
// terminated, all subsequent subscriptions to the observable side will be
// terminated immediately with either an Error or Complete notification send to
// the subscribing client
//
// Note that this implementation is blocking. When no subcribers are present
// then the data can flow freely. But when there are subscribers, the observable
// goroutine is blocked until all subscribers have processed the next, error or
// complete notification.
func NewSubjectFoo() SubjectFoo {
	ch := multicast.NewChan(1, 16 /*max endpoints*/)

	observable := Observable(func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		ep, err := ch.NewEndpoint(0)
		if err != nil {
			observe(nil, err, true)
			return
		}
		observable := Create(func(Next Next, Error Error, Complete Complete, Canceled Canceled) {
			receive := func(next interface{}, err error, closed bool) bool {
				switch {
				case !closed:
					Next(next)
				case err != nil:
					Error(err)
				default:
					Complete()
				}
				return !Canceled()
			}
			ep.Range(receive, 0)
		})
		observable(observe, subscribeOn, subscriber.Add(ep.Cancel))
	})

	observer := func(next foo, err error, done bool) {
		if !ch.Closed() {
			if !done {
				ch.FastSend(next)
			} else {
				ch.Close(err)
			}
		}
	}

	return SubjectFoo{observable.AsObservableFoo(), observer}
}

//jig:template MaxReplayCapacity

// MaxReplayCapacity is the maximum size of a replay buffer. Can be modified.
var MaxReplayCapacity = 16383

//jig:template NewReplaySubject<Foo>
//jig:needs Subject<Foo>, MaxReplayCapacity

// NewReplaySubjectFoo creates a new ReplaySubject. ReplaySubject ensures that
// all observers see the same sequence of emitted items, even if they
// subscribe after. When bufferCapacity argument is 0, then MaxReplayCapacity is
// used (currently 16383). When windowDuration argument is 0, then entries added
// to the buffer will remain fresh forever.
func NewReplaySubjectFoo(bufferCapacity int, windowDuration time.Duration) SubjectFoo {
	if bufferCapacity == 0 {
		bufferCapacity = MaxReplayCapacity
	}
	ch := multicast.NewChan(bufferCapacity, 16 /*max endpoints*/)

	observable := Observable(func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		ep, err := ch.NewEndpoint(multicast.ReplayAll)
		if err != nil {
			observe(nil, err, true)
			return
		}
		observable := Create(func(Next Next, Error Error, Complete Complete, Canceled Canceled) {
			receive := func(next interface{}, err error, closed bool) bool {
				switch {
				case !closed:
					Next(next)
				case err != nil:
					Error(err)
				default:
					Complete()
				}
				return !Canceled()
			}
			ep.Range(receive, windowDuration)
		})
		observable(observe, subscribeOn, subscriber.Add(ep.Cancel))
	})

	observer := func(next foo, err error, done bool) {
		if !ch.Closed() {
			if !done {
				ch.Send(next)
			} else {
				ch.Close(err)
			}
		}
	}

	return SubjectFoo{observable.AsObservableFoo(), observer}
}
