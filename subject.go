package rx

import (
	_ "github.com/reactivego/channel"
	"time"
)

//jig:template Subject<Foo>
//jig:embeds Observable<Foo>

// SubjectFoo is a combination of an observer and observable. Subjects are
// special because they are the only reactive constructs that support
// multicasting. The items sent to it through its observer side are
// multicasted to multiple clients subscribed to its observable side.
//
// A SubjectFoo embeds ObservableFoo and FooObserveFunc. This exposes the
// methods and fields of both types on SubjectFoo. Use the ObservableFoo
// methods to subscribe to it. Use the FooObserveFunc Next, Error and Complete
// methods to feed data to it.
//
// After a subject has been terminated by calling either Error or Complete,
// it goes into terminated state. All subsequent calls to its observer side
// will be silently ignored. All subsequent subscriptions to the observable
// side will be handled according to the specific behavior of the subject.
// There are different types of subjects, see the different NewXxxSubjectFoo
// functions for more info.
//
// Important! a subject is a hot observable. This means that subscribing to
// it will block the calling goroutine while it is waiting for items and
// notifications to receive. Unless you have code on a different goroutine
// already feeding into the subject, your subscribe will deadlock.
// Alternatively, you could subscribe on a goroutine as shown in the example.
type SubjectFoo struct {
	ObservableFoo
	FooObserveFunc
}

//jig:template NewSubject<Foo>
//jig:needs Subject<Foo>, NewChanFanOutNext

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
	channel := NewChanFanOutNext(1)

	observable := Observable(func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		channel, cancel := channel.NewChannel()
		observable := Create(func(observer Observer) {
			for item := range channel {
				if observer.Closed() {
					return
				}
				if item.Err == nil {
					observer.Next(item.Next)
				} else {
					observer.Error(item.Err)
					return
				}
			}
			observer.Complete()
		})
		observable(observe, subscribeOn, subscriber.Add(cancel))
	})

	observer := func(next foo, err error, done bool) {
		if !done {
			channel.Send(next)
		} else {
			channel.Close(err)
		}
	}

	return SubjectFoo{observable.AsFoo(), observer}
}

//jig:template MaxReplayCapacity

// MaxReplayCapacity is the maximum size of a replay buffer. Can be modified.
var MaxReplayCapacity = 16383

//jig:template NewReplaySubject<Foo>
//jig:needs Subject<Foo>, NewBufChan, MaxReplayCapacity, ErrTypecastTo<Foo>

// NewReplaySubjectFoo creates a new ReplaySubject. ReplaySubject ensures that
// all observers see the same sequence of emitted items, even if they
// subscribe after. When bufferCapacity argument is 0, then MaxReplayCapacity is
// used (currently 16383). When windowDuration argument is 0, then entries added
// to the buffer will remain fresh forever.
//
// Note that this implementation is non-blocking. When no subscribers are
// present the buffer fills up to bufferCapacity after which new items will
// start overwriting the oldest ones according to the FIFO principle.
// If a subscriber cannot keep up with the data rate of the source observable,
// eventually the buffer for the subscriber will overflow. At that moment the
// subscriber will receive an ErrMissingBackpressure error.
func NewReplaySubjectFoo(bufferCapacity int, windowDuration time.Duration) SubjectFoo {
	if bufferCapacity == 0 {
		bufferCapacity = MaxReplayCapacity
	}
	channel := NewBufChan(bufferCapacity, windowDuration)

	observable := Create(func(observer Observer) {
		channel.NewEndpoint().Range(func(value interface{}, err error, closed bool) bool {
			if observer.Closed() {
				return false
			}
			switch {
			case !closed:
				observer.Next(value)
			case err != nil:
				observer.Error(err)
			default:
				observer.Complete()
			}
			return !observer.Closed()
		})
	})

	observer := func(next foo, err error, done bool) {
		if !done {
			channel.Send(next)
		} else {
			channel.Close(err)
		}
	}

	return SubjectFoo{observable.AsFoo(), observer}
}
