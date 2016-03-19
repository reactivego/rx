package observable

import (
	"errors"
	"rxgo/examples/mousetest/mouse"
	"rxgo/observable/filters"
	"sync"
	"time"
)

// ErrTimeout is delivered to an observer if the stream times out.
var ErrTimeout = errors.New("timeout")

////////////////////////////////////////////////////////
// Scheduler
////////////////////////////////////////////////////////

type Scheduler interface {
	Schedule(task func())
}

type SchedulerFunc func(task func())

func (s SchedulerFunc) Schedule(task func()) {
	s(task)
}

var ImmediateScheduler SchedulerFunc = func(task func()) { task() }
var GoroutineScheduler SchedulerFunc = func(task func()) { go task() }

////////////////////////////////////////////////////////
// Executor
////////////////////////////////////////////////////////

type Executor interface {
	ExecuteOn(scheduler Scheduler)
}

type ExecutorFunc func(scheduler Scheduler)

func (e ExecutorFunc) ExecuteOn(scheduler Scheduler) {
	e(scheduler)
}

////////////////////////////////////////////////////////
// Int
////////////////////////////////////////////////////////

type Int func(IntObserverFunc) Unsubscriber

var zeroInt int

/////////////////////////////////////////////////////////////////////////////
// FROM
/////////////////////////////////////////////////////////////////////////////

type IntSubscriber interface {
	Next(int)
	Error(error)
	Complete()
	Unsubscribe()
	Unsubscribed() bool
}

type CreateIntFunc func(IntSubscriber)

func (f CreateIntFunc) Subscribe(observer IntObserverFunc) Unsubscriber {
	var subscriber struct {
		IntObserverFunc
		ExecutorFunc
		UnsubscriberInt
	}
	subscriber.IntObserverFunc = observer
	subscriber.ExecutorFunc = func(scheduler Scheduler) {
		task := func() {
			if !subscriber.Unsubscribed() {
				defer subscriber.Unsubscribe()
				f(&subscriber)
			}
		}
		scheduler.Schedule(task)
	}
	return &subscriber
}

// CreateInt calls f(subscriber) to produce values for a stream of ints.
func CreateInt(f CreateIntFunc) Int {
	return f.Subscribe
}

func Range(start, count int) Int {
	end := start + count
	return CreateInt(func(subscriber IntSubscriber) {
		for i := start; i < end; i++ {
			if subscriber.Unsubscribed() {
				return
			}
			subscriber.Next(i)
		}
		subscriber.Complete()
	})
}

func Interval(interval time.Duration) Int {
	return CreateInt(func(subscriber IntSubscriber) {
		i := 0
		for {
			time.Sleep(interval)
			if subscriber.Unsubscribed() {
				return
			}
			subscriber.Next(i)
			i++
		}
	})
}

// Repeat value count times.
func RepeatInt(value, count int) Int {
	return CreateInt(func(subscriber IntSubscriber) {
		for i := 0; i < count; i++ {
			if subscriber.Unsubscribed() {
				return
			}
			subscriber.Next(value)
		}
		subscriber.Complete()
	})
}

// StartInt is designed to be used with functions that return a
// (int, error) tuple.
//
// If the error is non-nil the returned Int will be that error,
// otherwise it will be a single-value stream of int.
func StartInt(f func() (int, error)) Int {
	return CreateInt(func(subscriber IntSubscriber) {
		if next, err := f(); err != nil {
			subscriber.Error(err)
		} else {
			subscriber.Next(next)
			subscriber.Complete()
		}
	})
}

func NeverInt() Int {
	return CreateInt(func(subscriber IntSubscriber) {
	})
}

func EmptyInt() Int {
	return CreateInt(func(subscriber IntSubscriber) {
		subscriber.Complete()
	})
}

func ThrowInt(err error) Int {
	return CreateInt(func(subscriber IntSubscriber) {
		subscriber.Error(err)
	})
}

func FromIntArray(array []int) Int {
	return CreateInt(func(subscriber IntSubscriber) {
		for _, next := range array {
			if subscriber.Unsubscribed() {
				return
			}
			subscriber.Next(next)
		}
		subscriber.Complete()
	})
}

func FromInts(array ...int) Int {
	return FromIntArray(array)
}

func JustInt(element int) Int {
	return CreateInt(func(subscriber IntSubscriber) {
		subscriber.Next(element)
		subscriber.Complete()
	})
}

func MergeInt(observables ...Int) Int {
	if len(observables) == 0 {
		return EmptyInt()
	}
	return (Int(observables[0].Subscribe)).Merge(observables[1:]...)
}

func MergeIntDelayError(observables ...Int) Int {
	if len(observables) == 0 {
		return EmptyInt()
	}
	return (Int(observables[0].Subscribe)).MergeDelayError(observables[1:]...)
}

func FromIntChannel(ch <-chan int) Int {
	return CreateInt(func(subscriber IntSubscriber) {
		for next := range ch {
			if subscriber.Unsubscribed() {
				return
			}
			subscriber.Next(next)
		}
		subscriber.Complete()
	})
}

/////////////////////////////////////////////////////////////////////////////
// Subscribe
/////////////////////////////////////////////////////////////////////////////

func (s Int) SubscribeOn(scheduler Scheduler) Int {
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		unsubscriber := s(observer)
		if executor, ok := unsubscriber.(Executor); ok {
			executor.ExecuteOn(scheduler)
		}
		return struct{ Unsubscriber }{unsubscriber}
	}
	return subscribe
}

func (s Int) Subscribe(observer IntObserverFunc) Unsubscriber {
	unsubscriber := s(observer)
	if executor, ok := unsubscriber.(Executor); ok {
		executor.ExecuteOn(GoroutineScheduler)
	}
	return struct{ Unsubscriber }{unsubscriber}
}

func (s Int) SubscribeNext(f func(v int)) Unsubscriber {
	return s.Subscribe(func(next int, err error, completed bool) {
		if err == nil && !completed {
			f(next)
		}
	})
}

// Wait for completion of the stream and return any error.
func (s Int) Wait() error {
	errch := make(chan error)
	s.Subscribe(func(next int, err error, completed bool) {
		if err != nil || completed {
			errch <- err
		}
	})
	return <-errch
}

/////////////////////////////////////////////////////////////////////////////
// TO
/////////////////////////////////////////////////////////////////////////////

// ToOneWithError blocks until the stream emits exactly one value. Otherwise, it errors.
func (s Int) ToOneWithError() (v int, e error) {
	v = zeroInt
	errch := make(chan error, 1)
	s.One().Subscribe(func(next int, err error, completed bool) {
		if err != nil || completed {
			errch <- err
			// Close errch to make subsequent use of it panic. This will prevent
			// a coroutine inside an observable getting stuck on errch and leaking.
			close(errch)
		} else {
			v = next
		}
	})
	e = <-errch
	return
}

// ToOne blocks and returns the only value emitted by the stream, or the zero
// value if an error occurs.
func (s Int) ToOne() int {
	value, _ := s.ToOneWithError()
	return value
}

// ToArrayWithError collects all values from the stream into an array,
// returning it and any error.
func (s Int) ToArrayWithError() (a []int, e error) {
	a = []int{}
	errch := make(chan error, 1)
	s.Subscribe(func(next int, err error, completed bool) {
		if err != nil || completed {
			errch <- err
			// Close errch to make subsequent use of it panic. This will prevent
			// a coroutine inside an observable getting stuck on errch and leaking.
			close(errch)
		} else {
			a = append(a, next)
		}
	})
	e = <-errch
	return
}

// ToArray blocks and returns the values from the stream in an array.
func (s Int) ToArray() []int {
	out, _ := s.ToArrayWithError()
	return out
}

// ToChannelWithError returns next and error channels corresponding to the stream elements and
// any error. When the next channel closes, it may be either because of an error or
// because the observable completed. The next channel may be closed without having emitted any
// values. The error channel always emits a value to indicate the observable has finished.
// When the error channel emits nil then the observable completed without errors, otherwise
// the error channel emits the error. When the observable has finished both channels will be
// closed.
func (s Int) ToChannelWithError() (<-chan int, <-chan error) {
	nextch := make(chan int, 1)
	errch := make(chan error, 1)
	s.Subscribe(func(next int, err error, completed bool) {
		if err != nil || completed {
			errch <- err
			close(errch)
			close(nextch)
		} else {
			nextch <- next
		}
	})
	return nextch, errch
}

func (s Int) ToChannel() <-chan int {
	ch, _ := s.ToChannelWithError()
	return ch
}

/////////////////////////////////////////////////////////////////////////////
// FILTERS
/////////////////////////////////////////////////////////////////////////////

func (s Int) WrapFilter(filter filters.Filter) Int {
	subscribe := func(sink IntObserverFunc) Unsubscriber {
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
		return s(intToGenericSource)
	}
	return subscribe
}

// Distinct removes duplicate elements in the stream.
func (s Int) Distinct() Int {
	return s.WrapFilter(filters.Distinct())
}

// ElementAt yields the Nth element of the stream.
func (s Int) ElementAt(n int) Int {
	return s.WrapFilter(filters.ElementAt(n))
}

// Filter elements in the stream on a function.
func (s Int) Filter(f func(int) bool) Int {
	predicate := func(v interface{}) bool {
		return f(v.(int))
	}
	return s.WrapFilter(filters.Where(predicate))
}

// Last returns just the first element of the stream.
func (s Int) First() Int {
	return s.WrapFilter(filters.First())
}

// Last returns just the last element of the stream.
func (s Int) Last() Int {
	return s.WrapFilter(filters.Last())
}

// SkipLast skips the first N elements of the stream.
func (s Int) Skip(n int) Int {
	return s.WrapFilter(filters.Skip(n))
}

// SkipLast skips the last N elements of the stream.
func (s Int) SkipLast(n int) Int {
	return s.WrapFilter(filters.SkipLast(n))
}

// Take returns just the first N elements of the stream.
func (s Int) Take(n int) Int {
	return s.WrapFilter(filters.Take(n))
}

// TakeLast returns just the last N elements of the stream.
func (s Int) TakeLast(n int) Int {
	return s.WrapFilter(filters.TakeLast(n))
}

// IgnoreElements ignores elements of the stream and emits only the completion events.
func (s Int) IgnoreElements() Int {
	return s.WrapFilter(filters.IgnoreElements())
}

// IgnoreCompletion ignores the completion event of the stream and therefore returns a stream that never completes.
func (s Int) IgnoreCompletion() Int {
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			if !completed {
				observer(next, err, completed)
			}
		}
		return s(operator)
	}
	return subscribe
}

func (s Int) One() Int {
	return s.WrapFilter(filters.One())
}

func (s Int) Replay(size int, duration time.Duration) Int {
	return s.WrapFilter(filters.Replay(size, duration))
}

func (s Int) Sample(duration time.Duration) Int {
	return s.WrapFilter(filters.Sample(duration))
}

// Debounce reduces subsequent duplicates to single items during a certain duration
func (s Int) Debounce(duration time.Duration) Int {
	return s.WrapFilter(filters.Debounce(duration))
}

/////////////////////////////////////////////////////////////////////////////
// COUNT
/////////////////////////////////////////////////////////////////////////////

func (s Int) Count() Int {
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		count := 0
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(count, nil, false)
				observer(zeroInt, err, completed)
			}
			count++
		}
		return s(operator)
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// CONCAT
/////////////////////////////////////////////////////////////////////////////

func (s Int) Concat(observables ...Int) Int {
	allObservables := append([]Int{s}, observables...)
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		unsubscribers := &UnsubscriberCollection{}

		var index int
		var maxIndex = len(allObservables) - 1

		var wg sync.WaitGroup
		operator := func(next int, err error, completed bool) {
			switch {
			case err != nil:
				index = maxIndex
				observer.Error(err)
				wg.Done()
			case completed:
				if index == maxIndex {
					observer.Complete()
				}
				wg.Done()
			default:
				observer.Next(next)
			}
		}

		execute := func(scheduler Scheduler) {
			scheduler.Schedule(func() {
				for i, o := range allObservables {
					index = i
					wg.Add(1)
					if !unsubscribers.Set(o.Subscribe(operator)) {
						return
					}
					wg.Wait()
					if index == maxIndex {
						return
					}
				}
			})
		}

		return &struct {
			ExecutorFunc
			Unsubscriber
		}{execute, unsubscribers}
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// MERGE
/////////////////////////////////////////////////////////////////////////////

func (s Int) merge(observables []Int, delayError bool) Int {
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		unsubscribers := &UnsubscriberCollection{}

		var (
			olock     sync.Mutex
			merged    int
			count     = len(observables)
			lastError error
		)
		operator := func(next int, err error, completed bool) {
			// Only one observable can be forwarding at any one time.
			olock.Lock()
			defer olock.Unlock()

			if merged >= count {
				return
			}

			switch {
			case err != nil:
				if delayError {
					lastError = err
					merged++
				} else {
					observer.Error(err)
					merged = count
					// tell all still executing observables, we're no longer interested
					unsubscribers.Unsubscribe()
				}

			case completed:
				merged++
				if merged == count {
					if lastError != nil {
						observer.Error(lastError)
					} else {
						observer.Complete()
					}
				}
			default:
				observer.Next(next)
			}
		}

		execute := func(scheduler Scheduler) {
			for _, subscribe := range observables {
				parent := subscribe(operator)
				if !unsubscribers.Add(parent) {
					return
				}
				if executor, ok := parent.(Executor); ok {
					executor.ExecuteOn(scheduler)
				}
			}
		}

		return &struct {
			ExecutorFunc
			Unsubscriber
		}{execute, unsubscribers}
	}
	return subscribe
}

// Merge an arbitrary number of observables with this one.
// An error from any of the observables will terminate the merged stream.
func (s Int) Merge(other ...Int) Int {
	if len(other) == 0 {
		return s
	}
	return s.merge(append(other, s), false)
}

// Merge an arbitrary number of observables with this one.
// Any error will be deferred until all observables terminate.
func (s Int) MergeDelayError(other ...Int) Int {
	if len(other) == 0 {
		return s
	}
	return s.merge(append(other, s), true)
}

/////////////////////////////////////////////////////////////////////////////
// CATCH
/////////////////////////////////////////////////////////////////////////////

func (s Int) Catch(catch Int) Int {
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		unsubscribers := &UnsubscriberCollection{}

		var thrownError error
		operator := func(next int, err error, completed bool) {
			if err != nil {
				thrownError = err
			} else {
				observer(next, err, completed)
			}
		}

		execute := func(scheduler Scheduler) {
			scheduler.Schedule(func() {
				parent := s(operator)
				if !unsubscribers.Set(parent) {
					return
				}
				if executor, ok := parent.(Executor); ok {
					executor.ExecuteOn(ImmediateScheduler)
				}
				if thrownError != nil {
					// The catch observable needs to be controlled by our
					// subscriber's Unsubscribe call and therefore is set
					// in the unsubscribers.
					// W.r.t. scheduling, the catch observable may have
					// its own configuration.
					unsubscribers.Set(catch.Subscribe(observer))
				}
			})
		}

		return &struct {
			ExecutorFunc
			Unsubscriber
		}{execute, unsubscribers}
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// RETRY
/////////////////////////////////////////////////////////////////////////////

func (s Int) Retry() Int {
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		unsubscribers := &UnsubscriberCollection{}

		finished := false
		operator := func(next int, err error, completed bool) {
			switch {
			case err != nil:
				// ignore
			case completed:
				finished = true
				observer.Complete()
			default:
				observer.Next(next)
			}
		}

		execute := func(scheduler Scheduler) {
			scheduler.Schedule(func() {
				for !finished {
					parent := s(operator)
					if !unsubscribers.Set(parent) {
						return
					}
					if executor, ok := parent.(Executor); ok {
						executor.ExecuteOn(ImmediateScheduler)
					}
				}
			})
		}

		return &struct {
			ExecutorFunc
			Unsubscriber
		}{execute, unsubscribers}
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// DO
/////////////////////////////////////////////////////////////////////////////

// Do applies a function for each value passing through the stream.
func (s Int) Do(f func(next int)) Int {
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			if err == nil && !completed {
				f(next)
			}
			observer(next, err, completed)
		}
		return s(operator)
	}
	return subscribe
}

// DoOnError applies a function for any error on the stream.
func (s Int) DoOnError(f func(err error)) Int {
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			if err != nil {
				f(err)
			}
			observer(next, err, completed)
		}
		return s(operator)
	}
	return subscribe
}

// DoOnComplete applies a function when the stream completes.
func (s Int) DoOnComplete(f func()) Int {
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			if completed {
				f()
			}
			observer(next, err, completed)
		}
		return s(operator)
	}
	return subscribe
}

// Finally applies a function for any error or completion on the stream.
// This doesn't expose whether this was an error or a completion.
func (s Int) Finally(f func()) Int {
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				f()
			}
			observer(next, err, completed)
		}
		return s(operator)
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// REDUCE
/////////////////////////////////////////////////////////////////////////////

func (s Int) Reduce(initial int, reducer func(int, int) int) Int {
	value := initial
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(value, nil, false)
				observer(zeroInt, err, completed)
			} else {
				value = reducer(value, next)
			}
		}
		return s(operator)
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// SCAN
/////////////////////////////////////////////////////////////////////////////

func (s Int) Scan(initial int, f func(int, int) int) Int {
	value := initial
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(zeroInt, err, completed)
			} else {
				value = f(value, next)
				observer(value, nil, false)
			}
		}
		return s(operator)
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// TIMEOUT
/////////////////////////////////////////////////////////////////////////////

func (s Int) Timeout(timeout time.Duration) Int {
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		cancel := s(observer)
		unsubscriber := NewUnsubscriberChannel()
		go func() {
			select {
			case <-time.After(timeout):
				observer.Error(ErrTimeout)
				cancel.Unsubscribe()
				unsubscriber.Unsubscribe()
			case <-unsubscriber:
				cancel.Unsubscribe()
			}
		}()
		return unsubscriber
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// FORK
/////////////////////////////////////////////////////////////////////////////

// Fork replicates each event from the parent to every subscriber of the fork.
func (s Int) Fork() Int {
	var (
		lock      sync.Mutex
		observers []IntObserverFunc
	)

	operator := func(next int, err error, completed bool) {
		lock.Lock()
		defer lock.Unlock()
		for _, observer := range observers {
			if observer != nil {
				observer(next, err, completed)
			}
		}
	}
	s(operator)

	subscribe := func(observer IntObserverFunc) Unsubscriber {
		lock.Lock()
		defer lock.Unlock()
		index := len(observers)
		observers = append(observers, observer)
		unsubscriber := UnsubscriberFunc(func() {
			lock.Lock()
			defer lock.Unlock()
			observers[index] = nil
		})
		return &unsubscriber
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// PUBLISH
/////////////////////////////////////////////////////////////////////////////

type ConnectableInt struct {
	Int
	connect func() Unsubscriber
}

// Connect will subscribe to the parent observable to start receiving values.
// All values will then be passed on to the observers that subscribed to this
// connectable observable
func (s ConnectableInt) Connect() Unsubscriber {
	return s.connect()
}

// Publish creates a connectable observable that only starts emitting values
// after the Connect method is called on it.
func (s Int) Publish() ConnectableInt {
	var (
		lock      sync.Mutex
		observers []IntObserverFunc
	)

	// Connect our operator observer function to the source start receiving
	// values from the source and forward them to our subscribers.
	connect := func() Unsubscriber {
		operator := func(next int, err error, completed bool) {
			lock.Lock()
			defer lock.Unlock()
			for _, observer := range observers {
				if observer != nil {
					observer(next, err, completed)
				}
			}
		}
		return s(operator)
	}

	// This subscribe function is installed so it is called by the code we inherit from Int
	// We should add subscribing observers to the slice of observers that are already subscribed.
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		lock.Lock()
		defer lock.Unlock()
		index := len(observers)
		observers = append(observers, observer)
		unsubscriber := UnsubscriberFunc(func() {
			lock.Lock()
			defer lock.Unlock()
			observers[index] = nil
		})
		return &unsubscriber
	}
	return ConnectableInt{Int: subscribe, connect: connect}
}

/////////////////////////////////////////////////////////////////////////////
// MATHEMATICAL
/////////////////////////////////////////////////////////////////////////////

func (s Int) Average() Int {
	var sum int
	var count int
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(sum/count, nil, false)
				observer(zeroInt, err, completed)
			} else {
				sum += next
				count++
			}
		}
		return s(operator)
	}
	return subscribe
}

func (s Int) Sum() Int {
	var sum int
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(sum, nil, false)
				observer(zeroInt, err, completed)
			} else {
				sum += next
			}
		}
		return s(operator)
	}
	return subscribe
}

func (s Int) Min() Int {
	started := false
	var min int
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				if started {
					observer(min, nil, false)
				}
				observer(zeroInt, err, completed)
			} else {
				if started {
					if min > next {
						min = next
					}
				} else {
					min = next
					started = true
				}
			}
		}
		return s(operator)
	}
	return subscribe
}

func (s Int) Max() Int {
	started := false
	var max int
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				if started {
					observer(max, nil, false)
				}
				observer(zeroInt, err, completed)
			} else {
				if started {
					if max < next {
						max = next
					}
				} else {
					max = next
					started = true
				}
			}
		}
		return s(operator)
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// MAP and FLATMAP (int -> int)
/////////////////////////////////////////////////////////////////////////////

func (s Int) Map(f func(int) int) Int {
	subscribe := func(observer IntObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			var mapped int
			if err == nil && !completed {
				mapped = f(next)
			}
			observer(mapped, err, completed)
		}
		return s(operator)
	}
	return subscribe
}

func (s Int) FlatMap(f func(int) Int) Int {
	subscribe := func(observer IntObserverFunc) Unsubscriber {

		subscription := new(UnsubscriberInt)
		var wg sync.WaitGroup
		var mx sync.Mutex
		gatedObserver := func(next int, err error, completed bool) {
			mx.Lock()
			defer mx.Unlock()
			observer(next, err, completed)
		}

		operator := func(next int, err error, completed bool) {
			switch {
			case err != nil:
				wg.Wait()
				observer.Error(err)
			case completed:
				wg.Wait()
				observer.Complete()
			default:
				wg.Add(1)
				f(next).Finally(func() { wg.Done() }).IgnoreCompletion().Subscribe(gatedObserver)
			}
		}
		s(operator)

		return subscription
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// MAP and FLATMAP (int -> float64)
/////////////////////////////////////////////////////////////////////////////

func (s Int) MapFloat64(f func(int) float64) Float64 {
	subscribe := func(observer Float64ObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			var mapped float64
			if err == nil && !completed {
				mapped = f(next)
			}
			observer(mapped, err, completed)
		}
		return s(operator)
	}
	return subscribe
}

func (s Int) FlatMapFloat64(f func(int) Float64) Float64 {
	subscribe := func(observer Float64ObserverFunc) Unsubscriber {
		var wg sync.WaitGroup
		var mx sync.Mutex
		gatedObserver := func(next float64, err error, completed bool) {
			mx.Lock()
			defer mx.Unlock()
			observer(next, err, completed)
		}
		operator := func(next int, err error, completed bool) {
			switch {
			case err != nil:
				wg.Wait()
				observer.Error(err)
			case completed:
				wg.Wait()
				observer.Complete()
			default:
				wg.Add(1)
				next := f(next)
				next.Finally(func() { wg.Done() }).IgnoreCompletion().Subscribe(gatedObserver)
			}
		}
		s(operator)
		return new(UnsubscriberInt)
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// MAP and FLATMAP (int -> string)
/////////////////////////////////////////////////////////////////////////////

func (s Int) MapString(f func(int) string) String {
	subscribe := func(observer StringObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			var mapped string
			if err == nil && !completed {
				mapped = f(next)
			}
			observer(mapped, err, completed)
		}
		return s(operator)
	}
	return subscribe
}

func (s Int) FlatMapString(f func(int) String) String {
	subscribe := func(observer StringObserverFunc) Unsubscriber {
		var wg sync.WaitGroup
		var mx sync.Mutex
		gatedObserver := func(next string, err error, completed bool) {
			mx.Lock()
			defer mx.Unlock()
			observer(next, err, completed)
		}
		operator := func(next int, err error, completed bool) {
			switch {
			case err != nil:
				wg.Wait()
				observer.Error(err)
			case completed:
				wg.Wait()
				observer.Complete()
			default:
				wg.Add(1)
				next := f(next)
				next.Finally(func() { wg.Done() }).IgnoreCompletion().Subscribe(gatedObserver)
			}
		}
		s(operator)
		return new(UnsubscriberInt)
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// MAP and FLATMAP (int -> *mouse.Move)
/////////////////////////////////////////////////////////////////////////////

func (s Int) MapMove(f func(int) *mouse.Move) Move {
	subscribe := func(observer MoveObserverFunc) Unsubscriber {
		operator := func(next int, err error, completed bool) {
			var mapped *mouse.Move
			if err == nil && !completed {
				mapped = f(next)
			}
			observer(mapped, err, completed)
		}
		return s(operator)
	}
	return subscribe
}

func (s Int) FlatMapMove(f func(int) Move) Move {
	subscribe := func(observer MoveObserverFunc) Unsubscriber {
		var wg sync.WaitGroup
		var mx sync.Mutex
		gatedObserver := func(next *mouse.Move, err error, completed bool) {
			mx.Lock()
			defer mx.Unlock()
			observer(next, err, completed)
		}
		operator := func(next int, err error, completed bool) {
			switch {
			case err != nil:
				wg.Wait()
				observer.Error(err)
			case completed:
				wg.Wait()
				observer.Complete()
			default:
				wg.Add(1)
				next := f(next)
				next.Finally(func() { wg.Done() }).IgnoreCompletion().Subscribe(gatedObserver)
			}
		}
		s(operator)
		return new(UnsubscriberInt)
	}
	return subscribe
}

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

////////////////////////////////////////////////////////
// Float64
////////////////////////////////////////////////////////

type Float64 func(Float64ObserverFunc) Unsubscriber

var zeroFloat64 float64

func (s Float64) Subscribe(observer Float64ObserverFunc) Unsubscriber {
	unsubscriber := s(observer)
	if executor, ok := unsubscriber.(Executor); ok {
		executor.ExecuteOn(GoroutineScheduler)
	}
	return unsubscriber
}

// IgnoreCompletion ignores the completion event of the stream and therefore returns a stream that never completes.
func (s Float64) IgnoreCompletion() Float64 {
	subscribe := func(observer Float64ObserverFunc) Unsubscriber {
		operator := func(next float64, err error, completed bool) {
			if !completed {
				observer(next, err, completed)
			}
		}
		return s(operator)
	}
	return subscribe
}

// Finally applies a function for any error or completion on the stream, using err == nil to indicate completion.
func (s Float64) Finally(f func()) Float64 {
	subscribe := func(observer Float64ObserverFunc) Unsubscriber {
		operator := func(next float64, err error, completed bool) {
			if err != nil || completed {
				f()
			}
			observer(next, err, completed)
		}
		return s(operator)
	}
	return subscribe
}

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

////////////////////////////////////////////////////////
// String
////////////////////////////////////////////////////////

type String func(StringObserverFunc) Unsubscriber

var zeroString string

func (s String) Subscribe(observer StringObserverFunc) Unsubscriber {
	unsubscriber := s(observer)
	if executor, ok := unsubscriber.(Executor); ok {
		executor.ExecuteOn(GoroutineScheduler)
	}
	return unsubscriber
}

// ToArrayWithError collects all values from the stream into an array,
// returning it and any error.
func (s String) ToArrayWithError() (a []string, e error) {
	a = []string{}
	errch := make(chan error, 1)
	s.Subscribe(func(next string, err error, completed bool) {
		if err != nil || completed {
			errch <- err
			// Close errch to make subsequent use of it panic. This will prevent
			// a coroutine inside an observable getting stuck on errch and leaking.
			close(errch)
		} else {
			a = append(a, next)
		}
	})
	e = <-errch
	return
}

// ToArray blocks and returns the values from the stream in an array.
func (s String) ToArray() []string {
	out, _ := s.ToArrayWithError()
	return out
}

// IgnoreCompletion ignores the completion event of the stream and therefore returns a stream that never completes.
func (s String) IgnoreCompletion() String {
	subscribe := func(observer StringObserverFunc) Unsubscriber {
		operator := func(next string, err error, completed bool) {
			if !completed {
				observer(next, err, completed)
			}
		}
		return s(operator)
	}
	return subscribe
}

// Finally applies a function for any error or completion on the stream, using err == nil to indicate completion.
func (s String) Finally(f func()) String {
	subscribe := func(observer StringObserverFunc) Unsubscriber {
		operator := func(next string, err error, completed bool) {
			if err != nil || completed {
				f()
			}
			observer(next, err, completed)
		}
		return s(operator)
	}
	return subscribe
}

type StringObserverFunc func(string, error, bool)

func (f StringObserverFunc) Next(next string) {
	f(next, nil, false)
}

func (f StringObserverFunc) Error(err error) {
	f(zeroString, err, false)
}

func (f StringObserverFunc) Complete() {
	f(zeroString, nil, true)
}

////////////////////////////////////////////////////////
// Move
////////////////////////////////////////////////////////

type Move func(MoveObserverFunc) Unsubscriber

var zeroMove *mouse.Move

func (s Move) Subscribe(observer MoveObserverFunc) Unsubscriber {
	unsubscriber := s(observer)
	if executor, ok := unsubscriber.(Executor); ok {
		executor.ExecuteOn(GoroutineScheduler)
	}
	return unsubscriber
}

// IgnoreCompletion ignores the completion event of the stream and therefore returns a stream that never completes.
func (s Move) IgnoreCompletion() Move {
	subscribe := func(observer MoveObserverFunc) Unsubscriber {
		operator := func(next *mouse.Move, err error, completed bool) {
			if !completed {
				observer(next, err, completed)
			}
		}
		return s(operator)
	}
	return subscribe
}

// Finally applies a function for any error or completion on the stream, using err == nil to indicate completion.
func (s Move) Finally(f func()) Move {
	subscribe := func(observer MoveObserverFunc) Unsubscriber {
		operator := func(next *mouse.Move, err error, completed bool) {
			if err != nil || completed {
				f()
			}
			observer(next, err, completed)
		}
		return s(operator)
	}
	return subscribe
}

type MoveObserverFunc func(*mouse.Move, error, bool)

func (f MoveObserverFunc) Next(next *mouse.Move) {
	f(next, nil, false)
}

func (f MoveObserverFunc) Error(err error) {
	f(zeroMove, err, false)
}

func (f MoveObserverFunc) Complete() {
	f(zeroMove, nil, true)
}
