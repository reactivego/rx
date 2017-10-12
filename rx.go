// Code generated by jig; DO NOT EDIT.

//go:generate jig --regen

package rx

import (
	"sync"
	"sync/atomic"
)

//jig:name BarObserveFunc

// BarObserveFunc is essentially the observer, a function that gets called
// whenever the observable has something to report.
type BarObserveFunc func(bar, error, bool)

var zeroBar bar

// Next is called by an ObservableBar to emit the next bar value to the
// observer.
func (f BarObserveFunc) Next(next bar) {
	f(next, nil, false)
}

// Error is called by an ObservableBar to report an error to the observer.
func (f BarObserveFunc) Error(err error) {
	f(zeroBar, err, true)
}

// Complete is called by an ObservableBar to signal that no more data is
// forthcoming to the observer.
func (f BarObserveFunc) Complete() {
	f(zeroBar, nil, true)
}

//jig:name ObservableBar

// ObservableBar is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableBar func(BarObserveFunc, Scheduler, Subscriber)

//jig:name ObserveFunc

// ObserveFunc is essentially the observer, a function that gets called
// whenever the observable has something to report.
type ObserveFunc func(interface{}, error, bool)

var zero interface{}

// Next is called by an Observable to emit the next interface{} value to the
// observer.
func (f ObserveFunc) Next(next interface{}) {
	f(next, nil, false)
}

// Error is called by an Observable to report an error to the observer.
func (f ObserveFunc) Error(err error) {
	f(zero, err, true)
}

// Complete is called by an Observable to signal that no more data is
// forthcoming to the observer.
func (f ObserveFunc) Complete() {
	f(zero, nil, true)
}

//jig:name Observable

// Observable is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type Observable func(ObserveFunc, Scheduler, Subscriber)

//jig:name IntObserveFunc

// IntObserveFunc is essentially the observer, a function that gets called
// whenever the observable has something to report.
type IntObserveFunc func(int, error, bool)

var zeroInt int

// Next is called by an ObservableInt to emit the next int value to the
// observer.
func (f IntObserveFunc) Next(next int) {
	f(next, nil, false)
}

// Error is called by an ObservableInt to report an error to the observer.
func (f IntObserveFunc) Error(err error) {
	f(zeroInt, err, true)
}

// Complete is called by an ObservableInt to signal that no more data is
// forthcoming to the observer.
func (f IntObserveFunc) Complete() {
	f(zeroInt, nil, true)
}

//jig:name ObservableInt

// ObservableInt is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableInt func(IntObserveFunc, Scheduler, Subscriber)

//jig:name ObservableFooObserveFunc

// ObservableFooObserveFunc is essentially the observer, a function that gets called
// whenever the observable has something to report.
type ObservableFooObserveFunc func(ObservableFoo, error, bool)

var zeroObservableFoo ObservableFoo

// Next is called by an ObservableObservableFoo to emit the next ObservableFoo value to the
// observer.
func (f ObservableFooObserveFunc) Next(next ObservableFoo) {
	f(next, nil, false)
}

// Error is called by an ObservableObservableFoo to report an error to the observer.
func (f ObservableFooObserveFunc) Error(err error) {
	f(zeroObservableFoo, err, true)
}

// Complete is called by an ObservableObservableFoo to signal that no more data is
// forthcoming to the observer.
func (f ObservableFooObserveFunc) Complete() {
	f(zeroObservableFoo, nil, true)
}

//jig:name ObservableObservableFoo

// ObservableObservableFoo is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableObservableFoo func(ObservableFooObserveFunc, Scheduler, Subscriber)

//jig:name ObservableFooMapObservableBar

// MapObservableBar transforms the items emitted by an ObservableFoo by applying a
// function to each item.
func (o ObservableFoo) MapObservableBar(project func(foo) ObservableBar) ObservableObservableBar {
	observable := func(observe ObservableBarObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
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

//jig:name IntObserver

// IntObserver is the interface used with CreateInt when implementing a custom
// observable.
type IntObserver interface {
	// Next emits the next int value.
	Next(int)
	// Error signals an error condition.
	Error(error)
	// Complete signals that no more data is to be expected.
	Complete()
	// Closed returns true when the subscription has been canceled.
	Closed() bool
}

//jig:name CreateInt

// CreateInt creates an Observable from scratch by calling observer methods
// programmatically.
func CreateInt(f func(IntObserver)) ObservableInt {
	observable := func(observe IntObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		scheduler.Schedule(func() {
			if subscriber.Closed() {
				return
			}
			observer := func(next int, err error, done bool) {
				if !subscriber.Closed() {
					observe(next, err, done)
				}
			}
			type observer_subscriber struct {
				IntObserveFunc
				Subscriber
			}
			f(&observer_subscriber{observer, subscriber})
		})
	}
	return observable
}

//jig:name Observer

// Observer is the interface used with Create when implementing a custom
// observable.
type Observer interface {
	// Next emits the next interface{} value.
	Next(interface{})
	// Error signals an error condition.
	Error(error)
	// Complete signals that no more data is to be expected.
	Complete()
	// Closed returns true when the subscription has been canceled.
	Closed() bool
}

//jig:name Create

// Create creates an Observable from scratch by calling observer methods
// programmatically.
func Create(f func(Observer)) Observable {
	observable := func(observe ObserveFunc, scheduler Scheduler, subscriber Subscriber) {
		scheduler.Schedule(func() {
			if subscriber.Closed() {
				return
			}
			observer := func(next interface{}, err error, done bool) {
				if !subscriber.Closed() {
					observe(next, err, done)
				}
			}
			type observer_subscriber struct {
				ObserveFunc
				Subscriber
			}
			f(&observer_subscriber{observer, subscriber})
		})
	}
	return observable
}

//jig:name Empty

// Empty creates an Observable that emits no items but terminates normally.
func Empty() Observable {
	return Create(func(observer Observer) {
		observer.Complete()
	})
}

//jig:name Next

// Next contains either the next interface{} value (in .Next) or an error (in .Err).
// If Err is nil then Next must be valid. Next is meant to be used as the
// type of a channel allowing errors to be delivered in-band with the values.
type Next struct {
	Next	interface{}
	Err	error
}

//jig:name NewChanNext

// ChanNext is a "chan Next" with two additional capabilities. Firstly, it
// can be properly canceled from the receiver side by calling the Cancel method.
// And secondly, it can deliver an error in-band (as opposed to out-of-band) to
// the receiver because of how Next is defined.
type ChanNext struct {
	// Channel can be used directly e.g. in a range statement to read Next
	// items from the channel.
	Channel	chan Next

	cancel	chan struct{}
	closed	bool
}

// NewChanNext creates a new ChanNext with given buffer capacity and
// returns a pointer to it. The Send and Close methods are supposed to be used
// from the sending side by a single goroutine. The Channel field and the Cancel
// method are supposed to be used from the receiving side and may be used from
// different goroutines. Multiple goroutines reading from a single channel will
// fight for values though, because a channel does not multicast.
func NewChanNext(capacity int) *ChanNext {
	return &ChanNext{
		Channel:	make(chan Next, capacity),
		cancel:		make(chan struct{}),
	}
}

// Send is used from the sender side to send the next value to the channel. This
// call only returns after delivering the value. If the channel has been closed,
// the value is ignored and the call returns immediately.
func (c *ChanNext) Send(value interface{}) bool {
	if c.closed {
		return false
	}
	select {
	case <-c.cancel:
		c.closed = true
		close(c.Channel)
		return false
	case c.Channel <- Next{Next: value}:
		return true
	}
}

// Close is used from the sender side to deliver an error value before closing
// the channel. Pass nil to indicate a normal close. If the channel has already
// been closed then the call will return immediately.
func (c *ChanNext) Close(err error) bool {
	if c.closed {
		return false
	}
	c.closed = true
	if err != nil {
		select {
		case <-c.cancel:
			return false
		case c.Channel <- Next{Err: err}:
		}
	}
	close(c.Channel)
	return true
}

// Cancel can be called exactly once from the receiver side to indicate it no
// longer is intereseted in the data or completion status. This cancelation will
// be signaled to the sender. The sender will be correctly aborted if it is
// already blocked on a Send or Close call to deliver a value or error to the
// receiver.
func (c *ChanNext) Cancel() {
	close(c.cancel)
}

//jig:name NewChanFanOutNext

// ChanFanOutNext is a blocking fan-out implementation writing to multiple
// ChanNext channels that implement message buffering. Use a call to
// NewChannel to add a receiving channel. A newly created channel will start
// receiving any messages that are send subsequently into the fan-out channel.
//
// This channel works by fanning out the Send calls to a slice of ChanNext
// wrapped Go channels. When one of the channels is slow and its buffer size
// reaches its maximum capacity, then the Send into that channel will block.
// This provides so called blocking backpressure to the sender by blocking its
// goroutine. Effectively the slowest channel will dictate the throughput for
// the whole fan-out assembly.
type ChanFanOutNext struct {
	sync.Mutex
	capacity	int
	channel		[]*ChanNext
	closed		bool
	err		error
}

// NewChanFanOutNext creates a new ChanFanOutNext with the given buffer
// capacity to use for the channels in the fanout
func NewChanFanOutNext(bufferCapacity int) *ChanFanOutNext {
	return &ChanFanOutNext{capacity: bufferCapacity}
}

// NewChannel adds a new ChanNext channel to the fan-out and returns its
// Channel field (<- chan Next) along with a cancel function. It is critical
// that the cancel function is called to indicate you want to stop receiving
// data, as simply abandoning the channel will fill up its buffer and then block
// the whole fan-out assembly from further processing messages.
//
// When the channel has been closed by the sender by calling Close, then the
// cancel function does not have to be called (but doing so does not hurt).
// But never call the cancel function more than once, because that will panic on
// closing an internal cancel channel twice.
//
// It is perfectly fine to call NewChannel on a fan-out channel that was already
// closed. The channel returned will replay the error if present and is then
// closed immediately.
func (m *ChanFanOutNext) NewChannel() (<-chan Next, func()) {
	m.Lock()
	defer m.Unlock()
	if m.closed {
		channel := make(chan Next, 1)
		if m.err != nil {
			channel <- Next{Err: m.err}
		}
		close(channel)
		return channel, func() {}
	}
	ch := NewChanNext(m.capacity)
	m.channel = append(m.channel, ch)
	return ch.Channel, func() {
		ch.Cancel()
		m.Lock()
		defer m.Unlock()
		for i, c := range m.channel {
			if c == ch {
				m.channel = append(m.channel[:i], m.channel[i+1:]...)
				return
			}
		}
	}
}

// Send is used to multicast a value to multiple receiving channels.
func (m *ChanFanOutNext) Send(value interface{}) {
	m.Lock()
	for _, c := range m.channel {
		c.Send(value)
	}
	m.Unlock()
}

// Close is used to close all receiving channels in the fan-out assembly. If
// err is not nil, then the error is send to the channels first before they
// are closed.
func (m *ChanFanOutNext) Close(err error) {
	m.Lock()
	for _, c := range m.channel {
		c.Close(err)
	}
	m.channel = nil
	m.closed = true
	m.err = err
	m.Unlock()
}

//jig:name ObservableBarObserveFunc

// ObservableBarObserveFunc is essentially the observer, a function that gets called
// whenever the observable has something to report.
type ObservableBarObserveFunc func(ObservableBar, error, bool)

var zeroObservableBar ObservableBar

// Next is called by an ObservableObservableBar to emit the next ObservableBar value to the
// observer.
func (f ObservableBarObserveFunc) Next(next ObservableBar) {
	f(next, nil, false)
}

// Error is called by an ObservableObservableBar to report an error to the observer.
func (f ObservableBarObserveFunc) Error(err error) {
	f(zeroObservableBar, err, true)
}

// Complete is called by an ObservableObservableBar to signal that no more data is
// forthcoming to the observer.
func (f ObservableBarObserveFunc) Complete() {
	f(zeroObservableBar, nil, true)
}

//jig:name ObservableObservableBar

// ObservableObservableBar is essentially a subscribe function taking an observe
// function, scheduler and an subscriber.
type ObservableObservableBar func(ObservableBarObserveFunc, Scheduler, Subscriber)

//jig:name ObservableSubscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscriber.
func (o Observable) Subscribe(observe ObserveFunc, setters ...SubscribeOptionSetter) Subscriber {
	scheduler := NewTrampoline()
	setter := SubscribeOn(scheduler, setters...)
	options := NewSubscribeOptions(setter)
	subscriber := options.NewSubscriber()
	observer := func(next interface{}, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			observe(zero, err, true)
			subscriber.Unsubscribe()
		}
	}
	o(observer, options.SubscribeOn, subscriber)
	return subscriber
}

//jig:name ObservableSerialize

// Serialize forces an Observable to make serialized calls and to be
// well-behaved.
func (o Observable) Serialize() Observable {
	observable := func(observe ObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			mutex		sync.Mutex
			alreadyDone	bool
		)
		observer := func(next interface{}, err error, done bool) {
			mutex.Lock()
			if !alreadyDone {
				alreadyDone = done
				observe(next, err, done)
			}
			mutex.Unlock()
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableObservableBarMergeAll

// MergeAll flattens a higher order observable by merging the observables it emits.
func (o ObservableObservableBar) MergeAll() ObservableBar {
	observable := func(observe BarObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			mutex	sync.Mutex
			count	int32	= 1
		)
		observer := func(next bar, err error, done bool) {
			mutex.Lock()
			defer mutex.Unlock()
			if !done || err != nil {
				observe(next, err, done)
			} else {
				if atomic.AddInt32(&count, -1) == 0 {
					observe(zeroBar, nil, true)
				}
			}
		}
		merger := func(next ObservableBar, err error, done bool) {
			if !done {
				atomic.AddInt32(&count, 1)
				next(observer, subscribeOn, subscriber)
			} else {
				observer(zeroBar, err, true)
			}
		}
		o(merger, subscribeOn, subscriber)
	}
	return observable
}
