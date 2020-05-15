// Code generated by jig; DO NOT EDIT.

//go:generate jig

package rx

import (
	"sync"
	"sync/atomic"

	"github.com/reactivego/subscriber"
)

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
	var zero interface{}
	observable := func(observe Observer, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Subscribed() {
				observe(zero, nil, true)
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name ObservableFooAsObservable

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

//jig:name ObservableFooMapObservableBar

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

//jig:name Next

// Next can be called to emit the next value to the IntObserver.
type Next func(interface{})

//jig:name Create

// Create provides a way of creating an Observable from
// scratch by calling observer methods programmatically.
//
// The create function provided to Create will be called once
// to implement the observable. It is provided with a Next, Error,
// Complete and Canceled function that can be called by the code that
// implements the Observable.
func Create(create func(Next, Error, Complete, Canceled)) Observable {
	var zero interface{}
	observable := func(observe Observer, scheduler Scheduler, subscriber Subscriber) {
		runner := scheduler.Schedule(func() {
			if subscriber.Canceled() {
				return
			}
			n := func(next interface{}) {
				if subscriber.Subscribed() {
					observe(next, nil, false)
				}
			}
			e := func(err error) {
				if subscriber.Subscribed() {
					observe(zero, err, true)
				}
			}
			c := func() {
				if subscriber.Subscribed() {
					observe(zero, nil, true)
				}
			}
			x := func() bool {
				return subscriber.Canceled()
			}
			create(n, e, c, x)
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:name FromObservableFoo

// FromObservableFoo creates an ObservableObservableFoo from multiple ObservableFoo values passed in.
func FromObservableFoo(slice ...ObservableFoo) ObservableObservableFoo {
	var zeroObservableFoo ObservableFoo
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
					observe(zeroObservableFoo, nil, true)
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

//jig:name ObservableSerialize

// Serialize forces an Observable to make serialized calls and to be
// well-behaved.
func (o Observable) Serialize() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var observer struct {
			sync.Mutex
			done	bool
		}
		serializer := func(next interface{}, err error, done bool) {
			observer.Lock()
			defer observer.Unlock()
			if !observer.done {
				observer.done = done
				observe(next, err, done)
			}
		}
		o(serializer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableBarAsObservable

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

//jig:name ObservableObservableBarCombineLatestAll

// CombineLatestAll flattens a higher order observable
// (e.g. ObservableObservableBar) by subscribing to
// all emitted observables (ie. ObservableBar entries) until the source
// completes. It will then wait for all of the subscribed ObservableBars
// to emit before emitting the first slice. Whenever any of the subscribed
// observables emits, a new slice will be emitted containing all the latest
// value.
func (o ObservableObservableBar) CombineLatestAll() ObservableBarSlice {
	var zeroBar bar
	var zeroBarSlice []bar
	observable := func(observe BarSliceObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observables := []ObservableBar(nil)
		var observers struct {
			sync.Mutex
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
						if observers.values[index] == zeroBar {
							observers.initialized++
						}
						observers.values[index] = next
						if observers.initialized == len(observers.values) {
							observe(observers.values, nil, false)
						}
					case err != nil:
						observers.active = 0
						observe(zeroBarSlice, err, true)
					default:
						if observers.active--; observers.active == 0 {
							observe(zeroBarSlice, nil, true)
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
				observe(zeroBarSlice, err, true)
			default:
				subscribeOn.Schedule(func() {
					if !subscriber.Canceled() {
						numObservables := len(observables)
						observers.values = make([]bar, numObservables)
						observers.active = numObservables
						for i, v := range observables {
							if subscriber.Canceled() {
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

//jig:name ObservableObservableBarConcatAll

// ConcatAll flattens a higher order observable by concattenating the observables it emits.
func (o ObservableObservableBar) ConcatAll() ObservableBar {
	var zeroBar bar
	observable := func(observe BarObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			mutex		sync.Mutex
			observables	[]ObservableBar
			observer	BarObserver
		)
		observer = func(next bar, err error, done bool) {
			mutex.Lock()
			defer mutex.Unlock()
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if len(observables) == 0 {
					observe(zeroBar, nil, true)
				} else {
					o := observables[0]
					observables = observables[1:]
					o(observer, subscribeOn, subscriber)
				}
			}
		}
		sourceSubscriber := subscriber.Add()
		concatenator := func(next ObservableBar, err error, done bool) {
			if !done {
				mutex.Lock()
				defer mutex.Unlock()
				observables = append(observables, next)
			} else {
				observer(zeroBar, err, done)
				sourceSubscriber.Unsubscribe()
			}
		}
		o(concatenator, subscribeOn, sourceSubscriber)
	}
	return observable
}

//jig:name ObservableObservableBarMergeAll

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
					var zeroBar bar
					observe(zeroBar, err, true)
				default:
					if atomic.AddInt32(&observers.len, -1) == 0 {
						var zeroBar bar
						observe(zeroBar, nil, true)
					}
				}
			}
		}
		merger := func(next ObservableBar, err error, done bool) {
			if !done {
				atomic.AddInt32(&observers.len, 1)
				next(observer, subscribeOn, subscriber)
			} else {
				var zeroBar bar
				observer(zeroBar, err, true)
			}
		}
		subscribeOn.Schedule(func() {
			if !subscriber.Canceled() {
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
			return RxError("Already Done")
		}
		return RxError("Recursion Error")
	}
	o.observe(o, next, err, done)
	if done {
		if err != nil {
			if !atomic.CompareAndSwapInt32(&o.state, linkBusy, linkError) {
				return RxError("Internal Error: 'busy' -> 'error'")
			}
		} else {
			if !atomic.CompareAndSwapInt32(&o.state, linkBusy, linkCompleting) {
				return RxError("Internal Error: 'busy' -> 'completing'")
			}
		}
	} else {
		if !atomic.CompareAndSwapInt32(&o.state, linkBusy, linkIdle) {
			return RxError("Internal Error: 'busy' -> 'idle'")
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
		return RxError("Already Subscribed")
	}
	observer := func(next bar, err error, done bool) {
		o.Observe(next, err, done)
	}
	observable(observer, scheduler, o.subscriber)
	if !atomic.CompareAndSwapInt32(&o.state, linkSubscribing, linkIdle) {
		return RxError("Internal Error")
	}
	return nil
}

func (o *linkBar) Cancel(callback func()) error {
	if !atomic.CompareAndSwapInt32(&o.callbackState, callbackNil, settingCallback) {
		return RxError("Already Waiting")
	}
	o.callbackKind = linkCancelOrCompleted
	o.callback = callback
	if !atomic.CompareAndSwapInt32(&o.callbackState, settingCallback, callbackSet) {
		return RxError("Internal Error")
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
		return RxError("Already Waiting")
	}
	o.callbackKind = linkCallbackOnComplete
	o.callback = callback
	if !atomic.CompareAndSwapInt32(&o.callbackState, settingCallback, callbackSet) {
		return RxError("Internal Error")
	}
	if atomic.CompareAndSwapInt32(&o.state, linkCompleting, linkComplete) {
		o.callback()
	}
	return nil
}

//jig:name ObservableObservableBarSwitchAll

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
					var zeroBar bar
					observe(zeroBar, err, true)
				})
				switcherSubscriber.Unsubscribe()
			default:
				currentLink.OnComplete(func() {
					var zeroBar bar
					observe(zeroBar, nil, true)
				})
				switcherSubscriber.Unsubscribe()
			}
		}
		o(switcher, subscribeOn, switcherSubscriber)
	}
	return observable
}
