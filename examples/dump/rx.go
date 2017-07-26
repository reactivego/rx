// This file was created by simulating what would happen if only code
// would be added that would fix compilation errors.
// To test the whole idea of just-in-time generic code expansion.

package main

import (
	"errors"
	"rxgo/filters"
	"rxgo/scheduler"
	"rxgo/unsubscriber"
	//"sync"
	//"time"
)

// ErrTimeout is delivered to an observer if the stream times out.
var ErrTimeout = errors.New("timeout")

// MaxReplaySize is the maximum size of a replay buffer. Can be modified.
var MaxReplaySize = 16384

type Scheduler scheduler.Scheduler

type Unsubscriber unsubscriber.Unsubscriber

////////////////////////////////////////////////////////
// IntObserver
////////////////////////////////////////////////////////

type IntObserver interface {
	Next(int)
	Error(error)
	Complete()
	Unsubscribed() bool
}

type IntObserverFunc func(int, error, bool)

var zeroInt int

func (f IntObserverFunc) Next(next int) {
	f(next, nil, false)
}

func (f IntObserverFunc) Error(err error) {
	f(zeroInt, err, err == nil)
}

func (f IntObserverFunc) Complete() {
	f(zeroInt, nil, true)
}

////////////////////////////////////////////////////////
// ObservableInt
////////////////////////////////////////////////////////

// Every observable is essentially a function taking an
// Observer function, scheduler and an unsubscriber.
type ObservableInt func(IntObserverFunc, Scheduler, Unsubscriber)

func CreateInt(f func(IntObserver)) ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		scheduler.Schedule(func() {
			if !unsubscriber.Unsubscribed() {
				operation := func(next int, err error, completed bool) {
					if !unsubscriber.Unsubscribed() {
						observer(next, err, completed)
					}
				}
				observer := &struct {
					IntObserverFunc
					Unsubscriber
				}{operation, unsubscriber}
				f(observer)
			}
		})
	}
	return observable
}

func FromIntArray(array []int) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		for _, next := range array {
			if observer.Unsubscribed() {
				return
			}
			observer.Next(next)
		}
		observer.Complete()
	})
}

func FromInts(array ...int) ObservableInt {
	return FromIntArray(array)
}

func (o ObservableInt) adaptFilter(filter filters.Filter) ObservableInt {
	observable := func(sink IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		genericToIntSink := func(next interface{}, err error, completed bool) {
			if nextInt, ok := next.(int); ok {
				sink(nextInt, err, completed)
			} else {
				sink(zeroInt, err, completed)
			}
		}
		genericToGenericFilter := filter(genericToIntSink)
		intToGenericSource := func(next int, err error, completed bool) {
			genericToGenericFilter(next, err, completed)
		}
		o(intToGenericSource, scheduler, unsubscriber)
	}
	return observable
}

func (o ObservableInt) Last() ObservableInt {
	return o.adaptFilter(filters.Last())
}

type Float64Observer interface {
	Next(float64)
	Error(error)
	Complete()
	Unsubscribed() bool
}

type Float64ObserverFunc func(float64, error, bool)

var zeroFloat64 float64

func (f Float64ObserverFunc) Next(next float64) {
	f(next, nil, false)
}

func (f Float64ObserverFunc) Error(err error) {
	f(zeroFloat64, err, err == nil)
}

func (f Float64ObserverFunc) Complete() {
	f(zeroFloat64, nil, true)
}

type ObservableFloat64 func(Float64ObserverFunc, Scheduler, Unsubscriber)

func (o ObservableInt) MapFloat64(f func(int) float64) ObservableFloat64 {
	observable := func(observer Float64ObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		operator := func(next int, err error, completed bool) {
			var mapped float64
			if err == nil && !completed {
				mapped = f(next)
			}
			observer(mapped, err, completed)
		}
		o(operator, scheduler, unsubscriber)
	}
	return observable
}
