package Publish

import (
	"sync"

	"github.com/reactivego/rx/scheduler"
	"github.com/reactivego/subscriber"
)

//jig:name Support

type Scheduler scheduler.Scheduler

type Subscriber subscriber.Subscriber

//jig:name IntObserverFunc

type IntObserverFunc func(int, error, bool)

var zeroInt int

func (f IntObserverFunc) Next(next int) {
	f(next, nil, false)
}

func (f IntObserverFunc) Error(err error) {
	f(zeroInt, err, true)
}

func (f IntObserverFunc) Complete() {
	f(zeroInt, nil, true)
}

//jig:name ObservableInt

// ObservableInt is essentially a function taking an observer function, scheduler and an subscriber.
type ObservableInt func(IntObserverFunc, Scheduler, Subscriber)

//jig:name CreateInt

// IntObserver is the interface used with CreateInt when implementing a custom observable.
type IntObserver interface {
	Next(int)
	Error(error)
	Complete()
	Closed() bool
}

// CreateInt creates an Observable from scratch by calling observer methods programmatically.
func CreateInt(f func(IntObserver)) ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, subscriber Subscriber) {
		scheduler.Schedule(func() {
			if !subscriber.Closed() {
				operation := func(next int, err error, done bool) {
					if !subscriber.Closed() {
						observer(next, err, done)
					}
				}
				observer := &struct {
					IntObserverFunc
					Subscriber
				}{operation, subscriber}
				f(observer)
			}
		})
	}
	return observable
}

//jig:name FromIntArray

// FromIntArray creates an Observable from a slice of values passed in.
func FromIntArray(array []int) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		for _, next := range array {
			if observer.Closed() {
				return
			}
			observer.Next(next)
		}
		observer.Complete()
	})
}

//jig:name FromInts

// FromInts creates an Observable from multiple values passed in as parameters.
func FromInts(array ...int) ObservableInt {
	return FromIntArray(array)
}

//jig:name ObservableIntSubscribeOn

// SubscribeOn specify the scheduler an Observable should use when it is subscribed to
func (o ObservableInt) SubscribeOn(subscribeOn Scheduler) ObservableInt {
	observable := func(observer IntObserverFunc, _ Scheduler, subscriber Subscriber) {
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ConnectableInt

// ConnectableInt is an ObservableInt that has an additional method Connect()
// used to Subscribe to the parent observable and then broadcasting values to
// all subscribers of ConnectableInt.
type ConnectableInt struct {
	ObservableInt
	connect	func() Subscriber
}

//jig:name ObservableIntPublish

// Publish converts an ordinary Observable into a connectable Observable.
// A connectable observable will only start emitting values after its Connect()
// method has been called.
func (o ObservableInt) Publish() ConnectableInt {

	type channel struct {
		data	chan int
		errs	chan error
	}

	var channels struct {
		sync.Mutex
		items	[]*channel
	}

	addChannel := func(data chan int, errs chan error) int {
		channels.Lock()
		defer channels.Unlock()
		index := len(channels.items)
		channels.items = append(channels.items, &channel{data, errs})
		return index
	}

	removeChannel := func(index int) {
		channels.Lock()
		defer channels.Unlock()
		ch := channels.items[index]
		if ch != nil {
			if ch.data != nil {
				close(ch.data)
			}
			channels.items[index] = nil
		}
	}

	removeAllChannels := func() {
		channels.Lock()
		defer channels.Unlock()
		for index, ch := range channels.items {
			if ch != nil {
				if ch.data != nil {
					close(ch.data)
				}
				channels.items[index] = nil
			}
		}
	}

	broadcastData := func(next int) {
		channels.Lock()
		defer channels.Unlock()
		for _, ch := range channels.items {
			if ch != nil {
				if ch.data != nil {
					ch.data <- next
				}
			}
		}
	}

	broadcastError := func(err error) {
		channels.Lock()
		defer channels.Unlock()
		for _, ch := range channels.items {
			if ch != nil {
				if ch.errs != nil {
					ch.errs <- err
				}
			}
		}
	}

	operator := func(next int, err error, done bool) {
		if !done {
			broadcastData(next)
		} else {
			if err != nil {
				broadcastError(err)
			} else {
				removeAllChannels()
			}
		}
	}

	connect := func() Subscriber {

		return o.Subscribe(operator)
	}

	observable := func(observer IntObserverFunc, subscribeOn Scheduler, subscriber Subscriber) {
		data := make(chan int, 1)
		errs := make(chan error, 1)
		index := addChannel(data, errs)

		subscriber.Add(func() {
			removeChannel(index)
		})
		FromIntChannelWithError(data, errs)(observer, subscribeOn, subscriber)
	}

	return ConnectableInt{ObservableInt: observable, connect: connect}
}

//jig:name ObservableIntMapString

// MapString transforms the items emitted by an ObservableInt by applying a function to each item.
func (o ObservableInt) MapString(project func(int) string) ObservableString {
	observable := func(observer StringObserverFunc, subscribeOn Scheduler, subscriber Subscriber) {
		operator := func(next int, err error, done bool) {
			var mapped string
			if !done {
				mapped = project(next)
			}
			observer(mapped, err, done)
		}
		o(operator, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableIntMapBool

// MapBool transforms the items emitted by an ObservableInt by applying a function to each item.
func (o ObservableInt) MapBool(project func(int) bool) ObservableBool {
	observable := func(observer BoolObserverFunc, subscribeOn Scheduler, subscriber Subscriber) {
		operator := func(next int, err error, done bool) {
			var mapped bool
			if !done {
				mapped = project(next)
			}
			observer(mapped, err, done)
		}
		o(operator, subscribeOn, subscriber)
	}
	return observable
}

//jig:name ObservableIntSubscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscriber.
func (o ObservableInt) Subscribe(observer IntObserverFunc) Subscriber {
	subscriber := subscriber.New()
	operator := func(next int, err error, done bool) {
		if done {
			subscriber.Unsubscribe()
		}
		observer(next, err, done)
	}
	o(operator, &scheduler.Trampoline{}, subscriber)
	return subscriber
}

//jig:name ConnectableIntConnect

// Connect instructs a connectable Observable to begin emitting items to its subscribers.
// All values will then be passed on to the observers that subscribed to this
// connectable observable
func (c ConnectableInt) Connect() Subscriber {
	return c.connect()
}

//jig:name FromIntChannelWithError

// FromIntChannelWithError creates an Observable from values and errors read from two separate go channels.
func FromIntChannelWithError(data <-chan int, errs <-chan error) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		for {
			select {
			case next, ok := <-data:
				if observer.Closed() {
					return
				}
				if ok {
					observer.Next(next)
				} else {
					data = nil
					observer.Complete()
					return
				}
			case err, ok := <-errs:
				if observer.Closed() {
					return
				}
				if ok {
					if err != nil {
						observer.Error(err)
						return
					}
				} else {
					errs = nil
				}
			}
		}
	})
}

//jig:name StringObserverFunc

type StringObserverFunc func(string, error, bool)

var zeroString string

func (f StringObserverFunc) Next(next string) {
	f(next, nil, false)
}

func (f StringObserverFunc) Error(err error) {
	f(zeroString, err, true)
}

func (f StringObserverFunc) Complete() {
	f(zeroString, nil, true)
}

//jig:name ObservableString

// ObservableString is essentially a function taking an observer function, scheduler and an subscriber.
type ObservableString func(StringObserverFunc, Scheduler, Subscriber)

//jig:name BoolObserverFunc

type BoolObserverFunc func(bool, error, bool)

var zeroBool bool

func (f BoolObserverFunc) Next(next bool) {
	f(next, nil, false)
}

func (f BoolObserverFunc) Error(err error) {
	f(zeroBool, err, true)
}

func (f BoolObserverFunc) Complete() {
	f(zeroBool, nil, true)
}

//jig:name ObservableBool

// ObservableBool is essentially a function taking an observer function, scheduler and an subscriber.
type ObservableBool func(BoolObserverFunc, Scheduler, Subscriber)

//jig:name ObservableStringSubscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscriber.
func (o ObservableString) Subscribe(observer StringObserverFunc) Subscriber {
	subscriber := subscriber.New()
	operator := func(next string, err error, done bool) {
		if done {
			subscriber.Unsubscribe()
		}
		observer(next, err, done)
	}
	o(operator, &scheduler.Trampoline{}, subscriber)
	return subscriber
}

//jig:name ObservableBoolSubscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscriber.
func (o ObservableBool) Subscribe(observer BoolObserverFunc) Subscriber {
	subscriber := subscriber.New()
	operator := func(next bool, err error, done bool) {
		if done {
			subscriber.Unsubscribe()
		}
		observer(next, err, done)
	}
	o(operator, &scheduler.Trampoline{}, subscriber)
	return subscriber
}
