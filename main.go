package main

import (
	"sync/atomic"
	"time"
)

////////////////////////////////////////////////////////
// main
////////////////////////////////////////////////////////

func main() {
	println("hello")

	// observable := Range(1, 10)
	observable := CreateInt(func(s IntSubscriber) {
		println("subscription received...")
		for i := 0; i < 15; i++ {
			for j := 0; j < 3; j++ {
				if s.Disposed() {
					return
				}
				s.Next(i)
			}
		}
		s.Complete()
		s.Dispose()
	})

	observable = observable.Distinct()

	term := make(chan struct{})
	subscription := observable.SubscribeFunc(func(v int, err error, complete bool) {
		if err != nil || complete {
			close(term)
			return
		}
		println(v)
	})
	<-term
	if subscription.Disposed() {
		println("subscription disposed....")
	} else {
		println("subscription still alive....")
	}
	println("goodbye")
}

////////////////////////////////////////////////////////
// Subscription
////////////////////////////////////////////////////////

type Subscription interface {
	Dispose()
	Disposed() bool
}

// Int32Subscription is an Int32 value with atomic operations implementing the Subscription interface

type Int32Subscription int32

func (t *Int32Subscription) Dispose() {
	atomic.StoreInt32((*int32)(t), 1)
}

func (t *Int32Subscription) Disposed() bool {
	return atomic.LoadInt32((*int32)(t)) == 1
}

////////////////////////////////////////////////////////
// Observable & Observer
////////////////////////////////////////////////////////

type Observer interface {
	Next(interface{})
	Error(error)
	Complete()
}

type Subscriber interface {
	Observer
	Subscription
}

type Observable interface {
	Subscribe(Observer) Subscription
}

type ObserverFunc func(interface{}, error, bool)

func (f ObserverFunc) Next(next interface{}) {
	f(next, nil, false)
}

func (f ObserverFunc) Error(err error) {
	f(nil, err, false)
}

func (f ObserverFunc) Complete() {
	f(nil, nil, true)
}

////////////////////////////////////////////////////////
// Int
////////////////////////////////////////////////////////

type IntObserver interface {
	Next(int)
	Error(error)
	Complete()
}

type IntSubscriber interface {
	IntObserver
	Subscription
}

type IntObservable interface {
	Subscribe(IntObserver) Subscription
}

var zeroInt = *new(int)

type IntObserverFunc func(int, error, bool)

func (f IntObserverFunc) Next(next int) {
	f(next, nil, false)
}

func (f IntObserverFunc) Error(err error) {
	f(zeroInt, err, false)
}

func (f IntObserverFunc) Complete() {
	f(zeroInt, nil, true)
}

type IntObservableFunc func(IntSubscriber)

func (f IntObservableFunc) Subscribe(observer IntObserver) Subscription {
	subscriber := &struct {
		IntObserver
		Int32Subscription
	}{IntObserver: observer}
	go f(subscriber)
	return subscriber
}

func PassthroughInt(next int, err error, complete bool, observer IntObserver) {
	switch {
	case err != nil:
		observer.Error(err)
	case complete:
		observer.Complete()
	default:
		observer.Next(next)
	}
}

type Int2IntFilter func(int, error, bool, IntObserver)
type Int2IntFilterFactory func(IntObserver) Int2IntFilter

type Int2IntStruct struct {
	parent        *IntStream
	filterFactory Int2IntFilterFactory
}

func (f *Int2IntStruct) Subscribe(observer IntObserver) Subscription {
	filter := f.filterFactory(observer)
	return f.parent.Subscribe(IntObserverFunc(func(next int, err error, complete bool) {
		filter(next, err, complete, observer)
	}))
}

func MakeInt2IntStreamWithParentAndFilter(parent *IntStream, filter Int2IntFilter) *IntStream {
	filterFactory := func(IntObserver) Int2IntFilter { return filter }
	return &IntStream{&Int2IntStruct{parent, filterFactory}}
}

func MakeInt2IntStreamWithParentAndMapper(parent *IntStream, mapper func(int) int) *IntStream {
	filterFactory := func(IntObserver) Int2IntFilter {
		filter := func(next int, err error, complete bool, observer IntObserver) {
			var mapped int
			if err == nil && !complete {
				mapped = mapper(next)
			}
			PassthroughInt(mapped, err, complete, observer)
		}
		return filter
	}
	return &IntStream{&Int2IntStruct{parent, filterFactory}}
}

func (ff FilterFactory) MakeInt2IntStreamWithParent(parent *IntStream) *IntStream {

	filterFactory := func(observer IntObserver) Int2IntFilter {
		downcastingSubscriber := &struct {
			ObserverFunc
			Int32Subscription
		}{
			ObserverFunc: func(next interface{}, err error, complete bool) {
				switch {
				case err != nil:
					observer.Error(err)
				case complete:
					observer.Complete()
				default:
					observer.Next(next.(int))
				}
			},
		}
		upcastingFilter := ff(downcastingSubscriber)
		return func(next int, err error, complete bool, observer IntObserver) {
			upcastingFilter(next, err, complete, downcastingSubscriber)
		}
	}

	return &IntStream{&Int2IntStruct{parent, filterFactory}}
}

/*
type MappingInt2Float64Func func(next int, err error, complete bool, observer Float64Observer)
type MappingInt2Float64FuncFactory func (observer Float64Observer) MappingInt2Float64Func

type MappingInt2Float64Observable struct {
	parent  IntObservable
	mapper MappingInt2Float64FuncFactory
}

func (f *MappingInt2Float64Observable) Subscribe(observer Float64Observer) Subscription {
	mapper := f.mapper(observer)
	return f.parent.Subscribe(IntObserverFunc(func(next int, err error, complete bool) {
		mapper(next, err, complete, observer)
	}))
}

func MapInt2Float64Observable(parent IntObservable, mapper MappingInt2Float64FuncFactory) Float64Observable {
	return &MappingInt2Float64Observable{
		parent:  parent,
		mapper: mapper,
	}
}

func MapInt2Float64ObserveDirect(parent IntObservable, mapper MappingInt2Float64Func) Float64Observable {
	return MapInt2Float64Observable(parent, func(Float64Observer) MappingInt2Float64Func {
		return mapper
	})
}

func MapInt2Float64ObserveNext(parent IntObservable, mapper func(int) float64) Float64Observable {
	return MapInt2Float64Observable(parent, func(Float64Observer) MappingInt2Float64Func {
			return func(next int, err error, complete bool, observer Float64Observer) {
				var mapped float64
				if err == nil && !complete {
					mapped = mapper(next)
				}
				PassthroughFloat64(mapped, err, complete, observer)
			}
		},
	)
}

type flatMapInt2Float64 struct {
	parent IntObservable
	mapper func (int) Float64Observable
}

func (f *flatMapInt2Float64) Subscribe(observer Float64Observer) Subscription {
	subscription := NewGenericSubscription()
	wg := sync.WaitGroup{}
	f.parent.Subscribe(IntObserverFunc(func (next int, err error, complete bool) {
		switch {
		case err != nil:
			wg.Wait()
			observer.Error(err)
		case complete:
			wg.Wait()
			observer.Complete()
		default:
			wg.Add(1)
			observable := f.mapper(next)
			stream := (&Float64Stream{observable}).
				DoOnComplete(func() { wg.Done() }).
				DoOnError(func(error) { wg.Done() })
			stream = &Float64Stream{ignoreCompletionFilter().Float64(stream)}
			stream.Subscribe(observer)
		}
	}))
	return subscription
}


// MapFloat64 maps this stream to an Float64Stream via f.
func (s *IntStream) MapFloat64(f func (int) float64) *Float64Stream {
	return FromFloat64Observable(MapInt2Float64ObserveNext(s, f))
}

func (s *IntStream) FlatMapFloat64(f func (int) Float64Observable) *Float64Stream {
	return &Float64Stream{&flatMapInt2Float64{s, f}}
}
*/

////////////////////////////
// IntStream
////////////////////////////

type IntStream struct {
	IntObservable
}

func CreateInt(f IntObservableFunc) *IntStream {
	return &IntStream{f}
}

func Range(start, count int) *IntStream {
	end := start + count
	return CreateInt(func(subscriber IntSubscriber) {
		for i := start; i < end; i++ {
			if subscriber.Disposed() {
				return
			}
			subscriber.Next(i)
		}
		subscriber.Complete()
		subscriber.Dispose()
	})
}

func (s *IntStream) SubscribeFunc(f IntObserverFunc) Subscription {
	return s.Subscribe(f)
}

// Distinct removes duplicate elements in the stream.
func (s *IntStream) Distinct() *IntStream {
	return distinctFilter().MakeInt2IntStreamWithParent(s)
}

// Debounce reduces subsequent duplicates to single items during a certain duration
func (s *IntStream) Debounce(duration time.Duration) *IntStream {
	return debounceFilter(duration).MakeInt2IntStreamWithParent(s)
}

// Wait for completion of the stream and return any error.
func (s *IntStream) Wait() error {
	errch := make(chan error)
	s.SubscribeFunc(func(next int, err error, complete bool) {
		switch {
		case err != nil:
			errch <- err
		case complete:
			errch <- nil
		default:
		}
	})
	return <-errch
}

func (s *IntStream) Do(f func(next int)) *IntStream {
	mapper := func(next int) int {
		f(next)
		return next
	}
	return MakeInt2IntStreamWithParentAndMapper(s, mapper)
}

func (s *IntStream) Reduce(initial int, reducer func(int, int) int) *IntStream {
	value := initial
	filter := func(next int, err error, complete bool, observer IntObserver) {
		switch {
		case err != nil:
			observer.Next(value)
			observer.Error(err)
		case complete:
			observer.Next(value)
			observer.Complete()
		default:
			value = reducer(value, next)
		}
	}
	return MakeInt2IntStreamWithParentAndFilter(s, filter)
}

func (s *IntStream) Scan(initial int, f func(int, int) int) *IntStream {
	value := initial
	filter := func(next int, err error, complete bool, observer IntObserver) {
		switch {
		case err != nil:
			observer.Error(err)
		case complete:
			observer.Complete()
		default:
			value = f(value, next)
			observer.Next(value)
		}
	}
	return MakeInt2IntStreamWithParentAndFilter(s, filter)
}

////////////////////////////////////
//  Generic Filter Implementations
////////////////////////////////////

type Filter func(interface{}, error, bool, Observer)

type FilterFactory func(Observer) Filter

func distinctFilter() FilterFactory {
	filterFactory := func(Observer) Filter {
		seen := map[interface{}]struct{}{}
		filter := func(next interface{}, err error, complete bool, observer Observer) {
			switch {
			case err != nil:
				observer.Error(err)
			case complete:
				observer.Complete()
			default:
				if _, ok := seen[next]; ok {
					return
				}
				seen[next] = struct{}{}
				observer.Next(next)
			}
		}
		return filter
	}
	return filterFactory
}

func debounceFilter(duration time.Duration) FilterFactory {
	filterFactory := func(observer Observer) Filter {
		errch := make(chan error)
		completech := make(chan bool)
		valuech := make(chan interface{})
		go func() {
			var timeout <-chan time.Time
			var nextValue interface{}
			for {
				select {
				case <-timeout:
					observer.Next(nextValue)
					timeout = nil
				case nextValue = <-valuech:
					timeout = time.After(duration)
				case err := <-errch:
					if timeout != nil {
						observer.Next(nextValue)
					}
					observer.Error(err)
					return
				case <-completech:
					if timeout != nil {
						observer.Next(nextValue)
					}
					observer.Complete()
					return
				}
			}
		}()
		filter := func(next interface{}, err error, complete bool, observer Observer) {
			switch {
			case err != nil:
				errch <- err
			case complete:
				completech <- true
			default:
				valuech <- next
			}
		}
		return filter
	}
	return filterFactory
}
