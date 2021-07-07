// Code generated by jig; DO NOT EDIT.

//go:generate jig

package rx

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/reactivego/subscriber"
)

//jig:name BarObserver

// BarObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type BarObserver func(next bar, err error, done bool)

//jig:name ObservableBar

// ObservableBar is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableBar func(BarObserver, Scheduler, Subscriber)

//jig:name Observer

// Observer is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type Observer func(next interface{}, err error, done bool)

//jig:name Observable

// Observable is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type Observable func(Observer, Scheduler, Subscriber)

//jig:name ObservableFooObserver

// ObservableFooObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type ObservableFooObserver func(next ObservableFoo, err error, done bool)

//jig:name ObservableObservableFoo

// ObservableObservableFoo is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableObservableFoo func(ObservableFooObserver, Scheduler, Subscriber)

//jig:name TimeObserver

// TimeObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type TimeObserver func(next time.Time, err error, done bool)

//jig:name ObservableTime

// ObservableTime is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableTime func(TimeObserver, Scheduler, Subscriber)

//jig:name TimestampFooObserver

// TimestampFooObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type TimestampFooObserver func(next TimestampFoo, err error, done bool)

//jig:name ObservableTimestampFoo

// ObservableTimestampFoo is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableTimestampFoo func(TimestampFooObserver, Scheduler, Subscriber)

//jig:name TimeIntervalFooObserver

// TimeIntervalFooObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type TimeIntervalFooObserver func(next TimeIntervalFoo, err error, done bool)

//jig:name ObservableTimeIntervalFoo

// ObservableTimeIntervalFoo is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableTimeIntervalFoo func(TimeIntervalFooObserver, Scheduler, Subscriber)

//jig:name BoolObserver

// BoolObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type BoolObserver func(next bool, err error, done bool)

//jig:name ObservableBool

// ObservableBool is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableBool func(BoolObserver, Scheduler, Subscriber)

//jig:name SliceObserver

// SliceObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type SliceObserver func(next Slice, err error, done bool)

//jig:name ObservableSlice

// ObservableSlice is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableSlice func(SliceObserver, Scheduler, Subscriber)

//jig:name FooSliceObserver

// FooSliceObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type FooSliceObserver func(next FooSlice, err error, done bool)

//jig:name ObservableFooSlice

// ObservableFooSlice is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableFooSlice func(FooSliceObserver, Scheduler, Subscriber)

//jig:name IntObserver

// IntObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type IntObserver func(next int, err error, done bool)

//jig:name ObservableInt

// ObservableInt is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableInt func(IntObserver, Scheduler, Subscriber)

//jig:name BarSliceObserver

// BarSliceObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type BarSliceObserver func(next BarSlice, err error, done bool)

//jig:name ObservableBarSlice

// ObservableBarSlice is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableBarSlice func(BarSliceObserver, Scheduler, Subscriber)

//jig:name Empty

// Empty creates an Observable that emits no items but terminates normally.
func Empty() Observable {
	observable := func(observe Observer, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Subscribed() {
				var zero interface{}
				observe(zero, nil, true)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name ObservableFoo_AsObservable

// AsObservable turns a typed ObservableFoo into an Observable of interface{}.
func (o ObservableFoo) AsObservable() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			observe(interface{}(next), err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableFoo_MapObservableBar

// MapObservableBar transforms the items emitted by an ObservableFoo by applying a
// function to each item.
func (o ObservableFoo) MapObservableBar(project func(foo) ObservableBar) ObservableObservableBar {
	observable := func(observe ObservableBarObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			var mapped ObservableBar
			if !done {
				mapped = project(next)
			}
			observe(mapped, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name Slice

type Slice = []interface{}

//jig:name Interval

// Interval creates an Observable that emits a sequence of integers spaced
// by a particular time interval. First integer is not emitted immediately, but
// only after the first time interval has passed. The generated code will do a type
// conversion from int to interface{}.
func Interval(interval time.Duration) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		runner := subscribeOn.ScheduleFutureRecursive(interval, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				observe(interface{}(i), nil, false)
				i++
				if subscriber.Subscribed() {
					self(interval)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name FromObservableFoo

// FromObservableFoo creates an ObservableObservableFoo from multiple ObservableFoo values passed in.
func FromObservableFoo(slice ...ObservableFoo) ObservableObservableFoo {
	observable := func(observe ObservableFooObserver, scheduler Scheduler, subscriber Subscriber) {
		i := 0
		runner := scheduler.ScheduleRecursive(func(self func()) {
			if subscriber.Subscribed() {
				if i < len(slice) {
					observe(slice[i], nil, false)
					if subscriber.Subscribed() {
						i++
						self()
					}
				} else {
					var zero ObservableFoo
					observe(zero, nil, true)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name BarSlice

type BarSlice = []bar

//jig:name ObservableBarObserver

// ObservableBarObserver is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type ObservableBarObserver func(next ObservableBar, err error, done bool)

//jig:name ObservableObservableBar

// ObservableObservableBar is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type ObservableObservableBar func(ObservableBarObserver, Scheduler, Subscriber)

//jig:name ObservableSlice_MapFooSlice

// MapFooSlice transforms the items emitted by an ObservableSlice by applying a
// function to each item.
func (o ObservableSlice) MapFooSlice(project func(Slice) FooSlice) ObservableFooSlice {
	observable := func(observe FooSliceObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next Slice, err error, done bool) {
			var mapped FooSlice
			if !done {
				mapped = project(next)
			}
			observe(mapped, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableBar_AsObservable

// AsObservable turns a typed ObservableBar into an Observable of interface{}.
func (o ObservableBar) AsObservable() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next bar, err error, done bool) {
			observe(interface{}(next), err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableObservableBar_CombineLatestAll

// CombineLatestAll flattens a higher order observable
// (e.g. ObservableObservableBar) by subscribing to
// all emitted observables (ie. ObservableBar entries) until the source
// completes. It will then wait for all of the subscribed ObservableBars
// to emit before emitting the first slice. Whenever any of the subscribed
// observables emits, a new slice will be emitted containing all the latest
// value.
func (o ObservableObservableBar) CombineLatestAll() ObservableBarSlice {
	observable := func(observe BarSliceObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observables := []ObservableBar(nil)
		var observers struct {
			sync.Mutex
			assigned	[]bool
			values		[]bar
			initialized	int
			active		int
		}
		makeObserver := func(index int) BarObserver {
			observer := func(next bar, err error, done bool) {
				observers.Lock()
				defer observers.Unlock()
				if observers.active > 0 {
					switch {
					case !done:
						if !observers.assigned[index] {
							observers.assigned[index] = true
							observers.initialized++
						}
						observers.values[index] = next
						if observers.initialized == len(observers.values) {
							observe(observers.values, nil, false)
						}
					case err != nil:
						observers.active = 0
						var zero []bar
						observe(zero, err, true)
					default:
						if observers.active--; observers.active == 0 {
							var zero []bar
							observe(zero, nil, true)
						}
					}
				}
			}
			return observer
		}

		observer := func(next ObservableBar, err error, done bool) {
			switch {
			case !done:
				observables = append(observables, next)
			case err != nil:
				var zero []bar
				observe(zero, err, true)
			default:
				subscribeOn.Schedule(func() {
					if subscriber.Subscribed() {
						numObservables := len(observables)
						observers.assigned = make([]bool, numObservables)
						observers.values = make([]bar, numObservables)
						observers.active = numObservables
						for i, v := range observables {
							if !subscriber.Subscribed() {
								return
							}
							v(makeObserver(i), subscribeOn, subscriber)
						}
					}
				})
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableObservableBar_ConcatAll

// ConcatAll flattens a higher order observable by concattenating the observables it emits.
func (o ObservableObservableBar) ConcatAll() ObservableBar {
	observable := func(observe BarObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var concat struct {
			sync.Mutex
			observables	[]ObservableBar
			observer	BarObserver
			subscriber	Subscriber
		}

		var source struct {
			observer	ObservableBarObserver
			subscriber	Subscriber
		}

		concat.observer = func(next bar, err error, done bool) {
			concat.Lock()
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if len(concat.observables) == 0 {
					if !source.subscriber.Subscribed() {
						var zero bar
						observe(zero, nil, true)
					}
					concat.observables = nil
				} else {
					observable := concat.observables[0]
					concat.observables = concat.observables[1:]
					observable(concat.observer, subscribeOn, subscriber)
				}
			}
			concat.Unlock()
		}

		source.observer = func(next ObservableBar, err error, done bool) {
			if !done {
				concat.Lock()
				initial := concat.observables == nil
				concat.observables = append(concat.observables, next)
				concat.Unlock()
				if initial {
					var zero bar
					concat.observer(zero, nil, true)
				}
			} else {
				concat.Lock()
				initial := concat.observables == nil
				source.subscriber.Done(err)
				concat.Unlock()
				if initial || err != nil {
					var zero bar
					concat.observer(zero, err, true)
				}
			}
		}
		source.subscriber = subscriber.Add()
		o(source.observer, subscribeOn, source.subscriber)
	}
	return observable
}

//jig:name ObservableObservableBar_MergeAll

// MergeAll flattens a higher order observable by merging the observables it emits.
func (o ObservableObservableBar) MergeAll() ObservableBar {
	observable := func(observe BarObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var observers struct {
			sync.Mutex
			done	bool
			len	int32
		}
		observer := func(next bar, err error, done bool) {
			observers.Lock()
			defer observers.Unlock()
			if !observers.done {
				switch {
				case !done:
					observe(next, nil, false)
				case err != nil:
					observers.done = true
					var zero bar
					observe(zero, err, true)
				default:
					if atomic.AddInt32(&observers.len, -1) == 0 {
						var zero bar
						observe(zero, nil, true)
					}
				}
			}
		}
		merger := func(next ObservableBar, err error, done bool) {
			if !done {
				atomic.AddInt32(&observers.len, 1)
				next(observer, subscribeOn, subscriber)
			} else {
				var zero bar
				observer(zero, err, true)
			}
		}
		subscribeOn.Schedule(func() {
			if subscriber.Subscribed() {
				observers.len = 1
				o(merger, subscribeOn, subscriber)
			}
		})
	}
	return observable
}

//jig:name linkBar

type linkBarObserver func(*linkBar, bar, error, bool)

type linkBar struct {
	observe		linkBarObserver
	state		int32
	callbackState	int32
	callbackKind	int
	callback	func()
	subscriber	Subscriber
}

func newInitialLinkBar() *linkBar {
	return &linkBar{state: linkCompleting, subscriber: subscriber.New()}
}

func newLinkBar(observe linkBarObserver, subscriber Subscriber) *linkBar {
	return &linkBar{
		observe:	observe,
		subscriber:	subscriber.Add(),
	}
}

func (o *linkBar) Observe(next bar, err error, done bool) error {
	if !atomic.CompareAndSwapInt32(&o.state, linkIdle, linkBusy) {
		if atomic.LoadInt32(&o.state) > linkBusy {
			return AlreadyDone
		}
		return RecursionNotAllowed
	}
	o.observe(o, next, err, done)
	if done {
		if err != nil {
			if !atomic.CompareAndSwapInt32(&o.state, linkBusy, linkError) {
				return StateTransitionFailed
			}
		} else {
			if !atomic.CompareAndSwapInt32(&o.state, linkBusy, linkCompleting) {
				return StateTransitionFailed
			}
		}
	} else {
		if !atomic.CompareAndSwapInt32(&o.state, linkBusy, linkIdle) {
			return StateTransitionFailed
		}
	}
	if atomic.LoadInt32(&o.callbackState) != callbackSet {
		return nil
	}
	if atomic.CompareAndSwapInt32(&o.state, linkCompleting, linkComplete) {
		o.callback()
	}
	if o.callbackKind == linkCancelOrCompleted {
		if atomic.CompareAndSwapInt32(&o.state, linkIdle, linkCanceled) {
			o.callback()
		}
	}
	return nil
}

func (o *linkBar) SubscribeTo(observable ObservableBar, scheduler Scheduler) error {
	if !atomic.CompareAndSwapInt32(&o.state, linkUnsubscribed, linkSubscribing) {
		return AlreadySubscribed
	}
	observer := func(next bar, err error, done bool) {
		o.Observe(next, err, done)
	}
	observable(observer, scheduler, o.subscriber)
	if !atomic.CompareAndSwapInt32(&o.state, linkSubscribing, linkIdle) {
		return StateTransitionFailed
	}
	return nil
}

func (o *linkBar) Cancel(callback func()) error {
	if !atomic.CompareAndSwapInt32(&o.callbackState, callbackNil, settingCallback) {
		return AlreadyWaiting
	}
	o.callbackKind = linkCancelOrCompleted
	o.callback = callback
	if !atomic.CompareAndSwapInt32(&o.callbackState, settingCallback, callbackSet) {
		return StateTransitionFailed
	}
	o.subscriber.Unsubscribe()
	if atomic.CompareAndSwapInt32(&o.state, linkCompleting, linkComplete) {
		o.callback()
	}
	if atomic.CompareAndSwapInt32(&o.state, linkIdle, linkCanceled) {
		o.callback()
	}
	return nil
}

func (o *linkBar) OnComplete(callback func()) error {
	if !atomic.CompareAndSwapInt32(&o.callbackState, callbackNil, settingCallback) {
		return AlreadyWaiting
	}
	o.callbackKind = linkCallbackOnComplete
	o.callback = callback
	if !atomic.CompareAndSwapInt32(&o.callbackState, settingCallback, callbackSet) {
		return StateTransitionFailed
	}
	if atomic.CompareAndSwapInt32(&o.state, linkCompleting, linkComplete) {
		o.callback()
	}
	return nil
}

//jig:name ObservableObservableBar_SwitchAll

// SwitchAll converts an Observable that emits Observables into a single Observable
// that emits the items emitted by the most-recently-emitted of those Observables.
func (o ObservableObservableBar) SwitchAll() ObservableBar {
	observable := func(observe BarObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(link *linkBar, next bar, err error, done bool) {
			if !done || err != nil {
				observe(next, err, done)
			} else {
				link.subscriber.Unsubscribe()
			}
		}
		currentLink := newInitialLinkBar()
		var switcherMutex sync.Mutex
		switcherSubscriber := subscriber.Add()
		switcher := func(next ObservableBar, err error, done bool) {
			switch {
			case !done:
				previousLink := currentLink
				func() {
					switcherMutex.Lock()
					defer switcherMutex.Unlock()
					currentLink = newLinkBar(observer, subscriber)
				}()
				previousLink.Cancel(func() {
					switcherMutex.Lock()
					defer switcherMutex.Unlock()
					currentLink.SubscribeTo(next, subscribeOn)
				})
			case err != nil:
				currentLink.Cancel(func() {
					var zero bar
					observe(zero, err, true)
				})
				switcherSubscriber.Unsubscribe()
			default:
				currentLink.OnComplete(func() {
					var zero bar
					observe(zero, nil, true)
				})
				switcherSubscriber.Unsubscribe()
			}
		}
		o(switcher, subscribeOn, switcherSubscriber)
	}
	return observable
}
