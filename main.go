package main

import (
	"sync"
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
//  Generic Filter Factory Implementations
////////////////////////////////////////////////////////

type FilterFactory func(ObserverFunc) ObserverFunc

type FiltersNamespace struct{}

var filters FiltersNamespace

func (FiltersNamespace) Passthrough() FilterFactory {
	factory := func(observer ObserverFunc) ObserverFunc {
		filter := func(next interface{}, err error, complete bool) {
			switch {
			case err != nil:
				observer.Error(err)
			case complete:
				observer.Complete()
			default:
				observer.Next(next)
			}
		}
		return filter
	}
	return factory
}
func (FiltersNamespace) IgnoreCompletion() FilterFactory {
	factory := func(observer ObserverFunc) ObserverFunc {
		filter := func(next interface{}, err error, complete bool) {
			switch {
			case err != nil:
				observer.Error(err)
			case complete:
				break
			default:
				observer.Next(next)
			}
		}
		return filter
	}
	return factory
}

func (FiltersNamespace) Distinct() FilterFactory {
	factory := func(observer ObserverFunc) ObserverFunc {
		seen := map[interface{}]struct{}{}
		filter := func(next interface{}, err error, complete bool) {
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
	return factory
}

func (FiltersNamespace) Debounce(duration time.Duration) FilterFactory {
	factory := func(observer ObserverFunc) ObserverFunc {
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
		filter := func(next interface{}, err error, complete bool) {
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
	return factory
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

///////// IntStream -> IntStream

type int2Int struct {
	source     *IntStream
	makeFilter func(IntObserver) IntObserverFunc
}

func (f *int2Int) Subscribe(observer IntObserver) Subscription {
	return f.source.Subscribe(f.makeFilter(observer))
}

func (makeGenericFilter FilterFactory) FilterIntStream(source *IntStream) *IntStream {

	makeFilter := func(sink IntObserver) IntObserverFunc {

		generic2Int := func(next interface{}, err error, complete bool) {
			switch {
			case err != nil:
				sink.Error(err)
			case complete:
				sink.Complete()
			default:
				sink.Next(next.(int))
			}
		}

		genericFilter := makeGenericFilter(generic2Int)

		int2Generic := func(next int, err error, complete bool) {
			genericFilter(next, err, complete)
		}
		return int2Generic
	}

	return &IntStream{&int2Int{source, makeFilter}}
}

type flatMapInt2Int struct {
	source IntObservable
	mapper func(int) IntObservable
}

func (f *flatMapInt2Int) Subscribe(observer IntObserver) Subscription {
	subscription := new(Int32Subscription)
	wg := sync.WaitGroup{}
	f.source.Subscribe(IntObserverFunc(func(next int, err error, complete bool) {
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
			stream := (&IntStream{observable}).DoOnComplete(wg.Done).DoOnError(func(error) { wg.Done() })
			stream = filters.IgnoreCompletion().FilterIntStream(stream)
			stream.Subscribe(observer)
		}
	}))
	return subscription
}

func (s *IntStream) FlatMap(f func(int) IntObservable) *IntStream {
	return &IntStream{&flatMapInt2Int{s, f}}
}

func (s *IntStream) Map(f func(int) int) *IntStream {
	factory := func(observer IntObserver) IntObserverFunc {
		filter := func(next int, err error, complete bool) {
			var mapped int
			if err == nil && !complete {
				mapped = f(next)
			}
			PassthroughInt(mapped, err, complete, observer)
		}
		return filter
	}
	return &IntStream{&int2Int{s, factory}}
}

// IntStream -> IntStream implementations

func (s *IntStream) Count() *IntStream {
	count := 0
	factory := func(observer IntObserver) IntObserverFunc {
		filter := func(next int, err error, complete bool) {
			switch {
			case err != nil:
				observer.Next(count)
				observer.Error(err)
			case complete:
				observer.Next(count)
				observer.Complete()
			default:
				count++
			}
		}
		return filter
	}
	return &IntStream{&int2Int{s, factory}}
}

// Distinct removes duplicate elements in the stream.
func (s *IntStream) Distinct() *IntStream {
	return filters.Distinct().FilterIntStream(s)
}

// Debounce reduces subsequent duplicates to single items during a certain duration
func (s *IntStream) Debounce(duration time.Duration) *IntStream {
	return filters.Debounce(duration).FilterIntStream(s)
}

// Do applies a function for each value passing through the stream.
func (s *IntStream) Do(f func(next int)) *IntStream {
	factory := func(observer IntObserver) IntObserverFunc {
		filter := func(next int, err error, complete bool) {
			if err == nil && !complete {
				f(next)
			}
			PassthroughInt(next, err, complete, observer)
		}
		return filter
	}
	return &IntStream{&int2Int{s, factory}}
}

// DoOnError applies a function for any error on the stream.
func (s *IntStream) DoOnError(f func(err error)) *IntStream {
	factory := func(observer IntObserver) IntObserverFunc {
		filter := func(next int, err error, complete bool) {
			if err != nil {
				f(err)
			}
			PassthroughInt(next, err, complete, observer)
		}
		return filter
	}
	return &IntStream{&int2Int{s, factory}}
}

// DoOnComplete applies a function when the stream completes.
func (s *IntStream) DoOnComplete(f func()) *IntStream {
	factory := func(observer IntObserver) IntObserverFunc {
		filter := func(next int, err error, complete bool) {
			if complete {
				f()
			}
			PassthroughInt(next, err, complete, observer)
		}
		return filter
	}
	return &IntStream{&int2Int{s, factory}}
}

func (s *IntStream) Reduce(initial int, reducer func(int, int) int) *IntStream {
	value := initial
	factory := func(observer IntObserver) IntObserverFunc {
		filter := func(next int, err error, complete bool) {
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
		return filter
	}
	return &IntStream{&int2Int{s, factory}}
}

func (s *IntStream) Scan(initial int, f func(int, int) int) *IntStream {
	value := initial
	factory := func(observer IntObserver) IntObserverFunc {
		filter := func(next int, err error, complete bool) {
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
		return filter
	}
	return &IntStream{&int2Int{s, factory}}
}

///////// IntStream -> Float64Stream

type Int2Float64 struct {
	source  IntObservable
	factory func(observer Float64Observer) IntObserverFunc
}

func (f *Int2Float64) Subscribe(observer Float64Observer) Subscription {
	filter := f.factory(observer)
	return f.source.Subscribe(IntObserverFunc(func(next int, err error, complete bool) {
		filter(next, err, complete)
	}))
}

func MakeInt2Float64StreamWithParentAndMapper(source *IntStream, mapper func(int) float64) *Float64Stream {
	factory := func(observer Float64Observer) IntObserverFunc {
		filter := func(next int, err error, complete bool) {
			var mapped float64
			if err == nil && !complete {
				mapped = mapper(next)
			}
			PassthroughFloat64(mapped, err, complete, observer)
		}
		return filter
	}
	return &Float64Stream{&Int2Float64{source, factory}}
}

/*
type flatMapInt2Float64 struct {
	source IntObservable
	mapper func(int) Float64Observable
}

func (f *flatMapInt2Float64) Subscribe(observer Float64Observer) Subscription {
	subscription := NewGenericSubscription()
	wg := sync.WaitGroup{}
	f.source.Subscribe(IntObserverFunc(func(next int, err error, complete bool) {
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
			stream = &Float64Stream{ignoreCompletion().FilterFloat64Stream(stream)}
			stream.Subscribe(observer)
		}
	}))
	return subscription
}

// IntStream -> Float64Stream implementations

func (s *IntStream) FlatMapFloat64(f func(int) Float64Observable) *Float64Stream {
	return &Float64Stream{&flatMapInt2Float64{s, f}}
}

// MapFloat64 maps this stream to an Float64Stream via f.
func (s *IntStream) MapFloat64(f func(int) float64) *Float64Stream {
	return FromFloat64Observable(MapInt2Float64ObserveNext(s, f))
}

*/

////////////////////////////////////////////////////////
// Float64
////////////////////////////////////////////////////////

type Float64Observer interface {
	Next(float64)
	Error(error)
	Complete()
}

type Float64Subscriber interface {
	Float64Observer
	Subscription
}

type Float64Observable interface {
	Subscribe(Float64Observer) Subscription
}

var zeroFloat64 = *new(float64)

type Float64ObserverFunc func(float64, error, bool)

func (f Float64ObserverFunc) Next(next float64) {
	f(next, nil, false)
}

func (f Float64ObserverFunc) Error(err error) {
	f(zeroFloat64, err, false)
}

func (f Float64ObserverFunc) Complete() {
	f(zeroFloat64, nil, true)
}

func PassthroughFloat64(next float64, err error, complete bool, observer Float64Observer) {
	switch {
	case err != nil:
		observer.Error(err)
	case complete:
		observer.Complete()
	default:
		observer.Next(next)
	}
}

type Float64Stream struct {
	Float64Observable
}

type float642Float64 struct {
	source  Float64Observable
	factory func(observer Float64Observer) Float64ObserverFunc
}

func (f *float642Float64) Subscribe(observer Float64Observer) Subscription {
	filter := f.factory(observer)
	return f.source.Subscribe(Float64ObserverFunc(func(next float64, err error, complete bool) {
		filter(next, err, complete)
	}))
}

func (ff FilterFactory) FilterFloat64Stream(source *Float64Stream) *Float64Stream {

	factory := func(observer Float64Observer) Float64ObserverFunc {
		downcastingSubscriber := func(next interface{}, err error, complete bool) {
			switch {
			case err != nil:
				observer.Error(err)
			case complete:
				observer.Complete()
			default:
				observer.Next(next.(float64))
			}
		}
		filter := ff(downcastingSubscriber)
		upcastingFilter := func(next float64, err error, complete bool) {
			filter(next, err, complete)
		}
		return upcastingFilter
	}

	return &Float64Stream{&float642Float64{source, factory}}
}
