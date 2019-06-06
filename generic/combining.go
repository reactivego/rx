package rx

import (
	"sync"
	"sync/atomic"

	"github.com/reactivego/subscriber"
)

//jig:template LinkEnums

// state
const (
	linkUnsubscribed = iota
	linkSubscribing
	linkIdle
	linkBusy
	linkError    // done:error
	linkCanceled // externally:canceled
	linkCompleting
	linkComplete // done:complete
)

// callbackState
const (
	callbackNil = iota
	settingCallback
	callbackSet
)

// callbackKind
const (
	linkCallbackOnComplete = iota
	linkCancelOrCompleted
)

//jig:template <Foo>Link
//jig:needs RxError, LinkEnums

type FooLinkObserveFunc func(*FooLink, foo, error, bool)

type FooLink struct {
	observe       FooLinkObserveFunc
	state         int32
	callbackState int32
	callbackKind  int
	callback      func()
	subscriber    Subscriber
}

func NewInitialFooLink() *FooLink {
	return &FooLink{state: linkCompleting, subscriber: subscriber.New()}
}

func NewFooLink(observe FooLinkObserveFunc, subscriber Subscriber) *FooLink {
	return &FooLink{
		observe:    observe,
		subscriber: subscriber.AddChild(),
	}
}

func (o *FooLink) Observe(next foo, err error, done bool) error {
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
		return nil // return when no close callback is set
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

func (o *FooLink) SubscribeTo(observable ObservableFoo, scheduler Scheduler) error {
	if !atomic.CompareAndSwapInt32(&o.state, linkUnsubscribed, linkSubscribing) {
		return RxError("Already Subscribed")
	}
	observer := func(next foo, err error, done bool) {
		o.Observe(next, err, done)
	}
	observable(observer, scheduler, o.subscriber)
	if !atomic.CompareAndSwapInt32(&o.state, linkSubscribing, linkIdle) {
		return RxError("Internal Error")
	}
	return nil
}

func (o *FooLink) Cancel(callback func()) error {
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

func (o *FooLink) OnComplete(callback func()) error {
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

//jig:template ObservableObservable<Foo> SwitchAll
//jig:needs <Foo>Link

// SwitchAll converts an Observable that emits Observables into a single Observable
// that emits the items emitted by the most-recently-emitted of those Observables.
func (o ObservableObservableFoo) SwitchAll() ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(link *FooLink, next foo, err error, done bool) {
			if !done || err != nil {
				observe(next, err, done)
			} else {
				link.subscriber.Unsubscribe() // We filter complete. Therefore, we need to perform Unsubscribe.
			}
		}
		currentLink := NewInitialFooLink()
		var switcherMutex sync.Mutex
		switcherSubscriber := subscriber.AddChild()
		switcher := func(next ObservableFoo, err error, done bool) {
			switch {
			case !done:
				previousLink := currentLink
				func() {
					switcherMutex.Lock()
					defer switcherMutex.Unlock()
					currentLink = NewFooLink(observer, subscriber)
				}()
				previousLink.Cancel(func() {
					switcherMutex.Lock()
					defer switcherMutex.Unlock()
					currentLink.SubscribeTo(next, subscribeOn)
				})
			case err != nil:
				currentLink.Cancel(func() {
					observe(zeroFoo, err, true)
				})
				switcherSubscriber.Unsubscribe()
			default:
				currentLink.OnComplete(func() {
					observe(zeroFoo, nil, true)
				})
				switcherSubscriber.Unsubscribe()
			}
		}
		o(switcher, subscribeOn, switcherSubscriber)
	}
	return observable
}

//jig:template Concat<Foo>
//jig:needs Observable<Foo> Concat

// ConcatFoo emits the emissions from two or more ObservableFoos without interleaving them.
func ConcatFoo(observables ...ObservableFoo) ObservableFoo {
	if len(observables) == 0 {
		return EmptyFoo()
	}
	return observables[0].Concat(observables[1:]...)
}

//jig:template Observable<Foo> Concat

// Concat emits the emissions from two or more ObservableFoos without interleaving them.
func (o ObservableFoo) Concat(other ...ObservableFoo) ObservableFoo {
	if len(other) == 0 {
		return o
	}
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			observables = append([]ObservableFoo{}, other...)
			observer    FooObserveFunc
		)
		observer = func(next foo, err error, done bool) {
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if len(observables) == 0 {
					observe(zeroFoo, nil, true)
				} else {
					o := observables[0]
					observables = observables[1:]
					o(observer, subscribeOn, subscriber)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template ObservableObservable<Foo> ConcatAll

// ConcatAll flattens a higher order observable by concattenating the observables it emits.
func (o ObservableObservableFoo) ConcatAll() ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			mutex       sync.Mutex
			observables []ObservableFoo
			observer    FooObserveFunc
		)
		observer = func(next foo, err error, done bool) {
			mutex.Lock()
			defer mutex.Unlock()
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if len(observables) == 0 {
					observe(zeroFoo, nil, true)
				} else {
					o := observables[0]
					observables = observables[1:]
					o(observer, subscribeOn, subscriber)
				}
			}
		}
		sourceSubscriber := subscriber.AddChild()
		concatenator := func(next ObservableFoo, err error, done bool) {
			if !done {
				mutex.Lock()
				defer mutex.Unlock()
				observables = append(observables, next)
			} else {
				observer(zeroFoo, err, done)
				sourceSubscriber.Unsubscribe()
			}
		}
		o(concatenator, subscribeOn, sourceSubscriber)
	}
	return observable
}

//jig:template Merge<Foo>
//jig:needs Observable<Foo> Merge

// MergeFoo combines multiple Observables into one by merging their emissions.
// An error from any of the observables will terminate the merged observables.
func MergeFoo(observables ...ObservableFoo) ObservableFoo {
	if len(observables) == 0 {
		return EmptyFoo()
	}
	return observables[0].Merge(observables[1:]...)
}

//jig:template Observable<Foo> Merge

// Merge combines multiple Observables into one by merging their emissions.
// An error from any of the observables will terminate the merged observables.
func (o ObservableFoo) Merge(other ...ObservableFoo) ObservableFoo {
	if len(other) == 0 {
		return o
	}
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			mutex sync.Mutex
			count = 1 + len(other)
		)
		observer := func(next foo, err error, done bool) {
			mutex.Lock()
			defer mutex.Unlock()
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if count--; count == 0 {
					observe(zeroFoo, nil, true)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
		for _, o := range other {
			o(observer, subscribeOn, subscriber)
		}
	}
	return observable
}

//jig:template ObservableObservable<Foo> MergeAll

// MergeAll flattens a higher order observable by merging the observables it emits.
func (o ObservableObservableFoo) MergeAll() ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			mutex sync.Mutex
			count int32 = 1
		)
		observer := func(next foo, err error, done bool) {
			mutex.Lock()
			defer mutex.Unlock()
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if atomic.AddInt32(&count, -1) == 0 {
					observe(zeroFoo, nil, true)
				}
			}
		}
		merger := func(next ObservableFoo, err error, done bool) {
			if !done {
				atomic.AddInt32(&count, 1)
				next(observer, subscribeOn, subscriber)
			} else {
				observer(zeroFoo, err, true)
			}
		}
		o(merger, subscribeOn, subscriber)
	}
	return observable
}

//jig:template MergeDelayError<Foo>
//jig:needs Observable<Foo> MergeDelayError

// MergeDelayErrorFoo combines multiple Observables into one by merging their emissions.
// Any error will be deferred until all observables terminate.
func MergeDelayErrorFoo(observables ...ObservableFoo) ObservableFoo {
	if len(observables) == 0 {
		return EmptyFoo()
	}
	return observables[0].MergeDelayError(observables[1:]...)
}

//jig:template Observable<Foo> MergeDelayError

// MergeDelayError combines multiple Observables into one by merging their emissions.
// Any error will be deferred until all observables terminate.
func (o ObservableFoo) MergeDelayError(other ...ObservableFoo) ObservableFoo {
	if len(other) == 0 {
		return o
	}
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			mutex      sync.Mutex
			count      = 1 + len(other)
			delayedErr error
		)
		observer := func(next foo, err error, done bool) {
			mutex.Lock()
			defer mutex.Unlock()
			if !done {
				observe(next, nil, false)
			} else {
				if err != nil {
					delayedErr = err
				}
				if count--; count == 0 {
					observe(zeroFoo, delayedErr, true)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
		for _, o := range other {
			o(observer, subscribeOn, subscriber)
		}
	}
	return observable
}
