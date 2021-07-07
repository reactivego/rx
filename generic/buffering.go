package rx

import (
	"sync"
	"time"
)

//jig:template Observable Buffer

// Buffer emits a sequence of Observables that emit what the source has emitted until the
// closingNotifier Observable emitted.
func (o Observable) Buffer(closingNotifier Observable) ObservableSlice {
	observable := func(observe SliceObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var serializer struct {
			sync.Mutex
			next []interface{}
			done bool
		}

		notifier := func(next interface{}, err error, done bool) {
			serializer.Lock()
			defer serializer.Unlock()
			if !serializer.done {
				serializer.done = done
				switch {
				case !done:
					observe(serializer.next, nil, false)
					serializer.next = serializer.next[:0]
				case err != nil:
					observe(nil, err, true)
					serializer.next = nil
				default:
					observe(serializer.next, nil, false)
					observe(nil, nil, true)
					serializer.next = nil
				}
			}
		}
		closingNotifier(notifier, subscribeOn, subscriber)

		observer := func(next interface{}, err error, done bool) {
			serializer.Lock()
			defer serializer.Unlock()
			if !serializer.done {
				serializer.done = done
				switch {
				case !done:
					serializer.next = append(serializer.next, next)
				case err != nil:
					observe(nil, err, true)
					serializer.next = nil
				default:
					observe(serializer.next, nil, false)
					observe(nil, nil, true)
					serializer.next = nil
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Buffer
//jig:needs Observable Buffer

// Buffer emits a sequence of Observables that emit what the source has emitted until the
// closingNotifier Observable emitted.
func (o ObservableFoo) Buffer(closingNotifier Observable) ObservableFooSlice {
	project := func(next Slice) FooSlice {
		foos := make(FooSlice, len(next))
		for i := range next {
			foos[i] = next[i].(foo)
		}
		return foos
	}
	return o.AsObservable().Buffer(closingNotifier).MapFooSlice(project)
}

//jig:template Observable<Foo> Buffer<Bar>
//jig:needs Observable<Foo> Buffer

// Buffer emits a sequence of Observables that emit what the source has emitted until the
// closingNotifier Observable emitted.
func (o ObservableFoo) BufferBar(closingNotifier ObservableBar) ObservableFooSlice {
	return o.Buffer(closingNotifier.AsObservable())
}

//jig:template Observable BufferTime

// BufferTime buffers the source Observable values for a specific time period
// and emits those as a slice periodically in time.
func (o Observable) BufferTime(period time.Duration) ObservableSlice {
	return o.Buffer(Interval(period))
}

//jig:template Observable<Foo> BufferTime
//jig:needs Observable BufferTime

// BufferTime buffers the source Observable values for a specific time period
// and emits those as a slice periodically in time.
func (o ObservableFoo) BufferTime(period time.Duration) ObservableFooSlice {
	return o.Buffer(Interval(period))
}
