package rx

import (
	"sync"
	"sync/atomic"

	"github.com/reactivego/rx/subscriber"
)

//jig:template Observable<Foo> SwitchMap<Bar>

// SwitchMapBar transforms the items emitted by an ObservableFoo by applying a
// function to each item an returning an ObservableBar. In doing so, it behaves much like
// what used to be called FlatMap, except that whenever a new ObservableBar is emitted
// SwitchMap will unsubscribe from the previous ObservableBar and begin emitting items
// from the newly emitted one.
func (o ObservableFoo) SwitchMapBar(project func(foo) ObservableBar) ObservableBar {
	return o.MapObservableBar(project).SwitchAll()
}

//jig:template LinkErrors

const (
	AlreadyDone           = RxError("already done")
	AlreadySubscribed     = RxError("already subscribed")
	AlreadyWaiting        = RxError("already waiting")
	RecursionNotAllowed   = RxError("recursion not allowed")
	StateTransitionFailed = RxError("state transition faled")
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

//jig:template link<Foo>
//jig:needs LinkErrors, LinkEnums

type linkFooObserver func(*linkFoo, foo, error, bool)

type linkFoo struct {
	observe       linkFooObserver
	state         int32
	callbackState int32
	callbackKind  int
	callback      func()
	subscriber    Subscriber
}

func newInitialLinkFoo() *linkFoo {
	return &linkFoo{state: linkCompleting, subscriber: subscriber.New()}
}

func newLinkFoo(observe linkFooObserver, subscriber Subscriber) *linkFoo {
	return &linkFoo{
		observe:    observe,
		subscriber: subscriber.Add(),
	}
}

func (o *linkFoo) Observe(next foo, err error, done bool) error {
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

func (o *linkFoo) SubscribeTo(observable ObservableFoo, scheduler Scheduler) error {
	if !atomic.CompareAndSwapInt32(&o.state, linkUnsubscribed, linkSubscribing) {
		return AlreadySubscribed
	}
	observer := func(next foo, err error, done bool) {
		o.Observe(next, err, done)
	}
	observable(observer, scheduler, o.subscriber)
	if !atomic.CompareAndSwapInt32(&o.state, linkSubscribing, linkIdle) {
		return StateTransitionFailed
	}
	return nil
}

func (o *linkFoo) Cancel(callback func()) error {
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

func (o *linkFoo) OnComplete(callback func()) error {
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

//jig:template ObservableObservable<Foo> SwitchAll
//jig:needs link<Foo>

// SwitchAll converts an Observable that emits Observables into a single Observable
// that emits the items emitted by the most-recently-emitted of those Observables.
func (o ObservableObservableFoo) SwitchAll() ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(link *linkFoo, next foo, err error, done bool) {
			if !done || err != nil {
				observe(next, err, done)
			} else {
				link.subscriber.Unsubscribe() // We filter complete. Therefore, we need to perform Unsubscribe.
			}
		}
		currentLink := newInitialLinkFoo()
		var switcherMutex sync.Mutex
		switcherSubscriber := subscriber.Add()
		switcher := func(next ObservableFoo, err error, done bool) {
			switch {
			case !done:
				previousLink := currentLink
				func() {
					switcherMutex.Lock()
					defer switcherMutex.Unlock()
					currentLink = newLinkFoo(observer, subscriber)
				}()
				previousLink.Cancel(func() {
					switcherMutex.Lock()
					defer switcherMutex.Unlock()
					currentLink.SubscribeTo(next, subscribeOn)
				})
			case err != nil:
				currentLink.Cancel(func() {
					var zero foo
					observe(zero, err, true)
				})
				switcherSubscriber.Unsubscribe()
			default:
				currentLink.OnComplete(func() {
					var zero foo
					observe(zero, nil, true)
				})
				switcherSubscriber.Unsubscribe()
			}
		}
		o(switcher, subscribeOn, switcherSubscriber)
	}
	return observable
}
