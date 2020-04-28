package rx

import (
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

//jig:template link<Foo>
//jig:needs RxError, LinkEnums

type linkFooObserveFunc func(*linkFoo, foo, error, bool)

type linkFoo struct {
	observe       linkFooObserveFunc
	state         int32
	callbackState int32
	callbackKind  int
	callback      func()
	subscriber    Subscriber
}

func newInitialLinkFoo() *linkFoo {
	return &linkFoo{state: linkCompleting, subscriber: subscriber.New()}
}

func newLinkFoo(observe linkFooObserveFunc, subscriber Subscriber) *linkFoo {
	return &linkFoo{
		observe:    observe,
		subscriber: subscriber.Add(),
	}
}

func (o *linkFoo) Observe(next foo, err error, done bool) error {
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

func (o *linkFoo) SubscribeTo(observable ObservableFoo, scheduler Scheduler) error {
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

func (o *linkFoo) Cancel(callback func()) error {
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

func (o *linkFoo) OnComplete(callback func()) error {
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
