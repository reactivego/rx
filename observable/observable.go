package observable

import (
	"errors"
	"rxgo/filters"
	"rxgo/schedulers"
	"rxgo/unsubscriber"
	"sync"
	"time"
)

// ErrTimeout is delivered to an observer if the stream times out.
var ErrTimeout = errors.New("timeout")

// MaxReplaySize is the maximum size of a replay buffer. Can be modified.
var MaxReplaySize = 16384

type Scheduler schedulers.Scheduler

type Unsubscriber unsubscriber.Unsubscriber

// Subscribable is an interface returned as a result of calling
// the observable function. See e.g. type ObservableInt. Every
// observable is essentially a function taking an Observer (function)
// as parameter and returning a Subscribable (interface).
type Subscribable interface {
	SubscribeOn(scheduler Scheduler)
	Unsubscriber
}

type SubscribeOnFunc func(scheduler Scheduler)

func (s SubscribeOnFunc) SubscribeOn(scheduler Scheduler) {
	s(scheduler)
}

////////////////////////////////////////////////////////
// ObservableInt
////////////////////////////////////////////////////////

type ObservableInt func(IntObserverFunc) Subscribable

type IntObserver interface {
	Next(int)
	Error(error)
	Complete()
}

type IntObserverFunc func(int, error, bool)

var zeroInt int

func (f IntObserverFunc) Next(next int) {
	f(next, nil, false)
}

func (f IntObserverFunc) Error(err error) {
	f(zeroInt, err, false)
}

func (f IntObserverFunc) Complete() {
	f(zeroInt, nil, true)
}

/////////////////////////////////////////////////////////////////////////////
// FROM
/////////////////////////////////////////////////////////////////////////////

type IntSubscriber interface {
	IntObserver
	Unsubscriber
}

type IntCreateFunc func(IntSubscriber)

// CreateInt calls f(subscriber) to produce values for a stream of ints.
func CreateInt(f IntCreateFunc) ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
		unsubscriber := new(unsubscriber.Int32)

		operation := func(next int, err error, completed bool) {
			if !unsubscriber.Unsubscribed() {
				observer(next, err, completed)
			}
		}

		subscribeOn := func(scheduler Scheduler) {
			scheduler.Schedule(func() {
				if !unsubscriber.Unsubscribed() {
					defer unsubscriber.Unsubscribe()
					subscriber := &struct {
						IntObserverFunc
						Unsubscriber
					}{operation, unsubscriber}
					f(subscriber)
				}
			})
		}

		return &struct {
			SubscribeOnFunc
			Unsubscriber
		}{subscribeOn, unsubscriber}
	}
	return subscribe
}

func EmptyInt() ObservableInt {
	return CreateInt(func(subscriber IntSubscriber) {
		subscriber.Complete()
	})
}

func NeverInt() ObservableInt {
	return CreateInt(func(subscriber IntSubscriber) {
	})
}

func ThrowInt(err error) ObservableInt {
	return CreateInt(func(subscriber IntSubscriber) {
		subscriber.Error(err)
	})
}

func FromIntArray(array []int) ObservableInt {
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

func FromInts(array ...int) ObservableInt {
	return FromIntArray(array)
}

func FromIntChannel(ch <-chan int) ObservableInt {
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

func Interval(interval time.Duration) ObservableInt {
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

func JustInt(element int) ObservableInt {
	return CreateInt(func(subscriber IntSubscriber) {
		subscriber.Next(element)
		subscriber.Complete()
	})
}

func Range(start, count int) ObservableInt {
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

// Repeat value count times.
func RepeatInt(value, count int) ObservableInt {
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
// If the error is non-nil the returned ObservableInt will be that error,
// otherwise it will be a single-value stream of int.
func StartInt(f func() (int, error)) ObservableInt {
	return CreateInt(func(subscriber IntSubscriber) {
		if next, err := f(); err != nil {
			subscriber.Error(err)
		} else {
			subscriber.Next(next)
			subscriber.Complete()
		}
	})
}

func MergeInt(observables ...ObservableInt) ObservableInt {
	if len(observables) == 0 {
		return EmptyInt()
	}
	return observables[0].Merge(observables[1:]...)
}

func MergeIntDelayError(observables ...ObservableInt) ObservableInt {
	if len(observables) == 0 {
		return EmptyInt()
	}
	return observables[0].MergeDelayError(observables[1:]...)
}

/////////////////////////////////////////////////////////////////////////////
// Subscribe
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) SubscribeOn(scheduler Scheduler) ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
		subscribable := s(observer)
		subscribable.SubscribeOn(scheduler)
		subscribeOn := func(scheduler Scheduler) {
			// ignore
		}
		return &struct {
			SubscribeOnFunc
			Unsubscriber
		}{subscribeOn, subscribable}
	}
	return subscribe
}

func (s ObservableInt) Subscribe(observer IntObserverFunc) Unsubscriber {
	subscribable := s(observer)
	subscribable.SubscribeOn(schedulers.GoroutineScheduler)
	// Only export Unsubscriber interface, not whole Subscribable.
	return &struct{ Unsubscriber }{subscribable}
}

func (s ObservableInt) SubscribeNext(f func(v int)) Unsubscriber {
	operator := func(next int, err error, completed bool) {
		if err == nil && !completed {
			f(next)
		}
	}
	return s.Subscribe(operator)
}

// Wait for completion of the stream and return any error.
func (s ObservableInt) Wait() error {
	doneChan := make(chan error)
	operator := func(next int, err error, completed bool) {
		if err != nil || completed {
			doneChan <- err
		}
	}
	s.Subscribe(operator)
	return <-doneChan
}

/////////////////////////////////////////////////////////////////////////////
// TO
/////////////////////////////////////////////////////////////////////////////

// ToOneWithError blocks until the stream emits exactly one value. Otherwise, it errors.
func (s ObservableInt) ToOneWithError() (v int, e error) {
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
func (s ObservableInt) ToOne() int {
	value, _ := s.ToOneWithError()
	return value
}

// ToArrayWithError collects all values from the stream into an array,
// returning it and any error.
func (s ObservableInt) ToArrayWithError() (a []int, e error) {
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
func (s ObservableInt) ToArray() []int {
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
func (s ObservableInt) ToChannelWithError() (<-chan int, <-chan error) {
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

func (s ObservableInt) ToChannel() <-chan int {
	ch, _ := s.ToChannelWithError()
	return ch
}

/////////////////////////////////////////////////////////////////////////////
// FILTERS
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) adaptFilter(filter filters.Filter) ObservableInt {
	subscribe := func(sink IntObserverFunc) Subscribable {
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
func (s ObservableInt) Distinct() ObservableInt {
	return s.adaptFilter(filters.Distinct())
}

// ElementAt yields the Nth element of the stream.
func (s ObservableInt) ElementAt(n int) ObservableInt {
	return s.adaptFilter(filters.ElementAt(n))
}

// Filter elements in the stream on a function.
func (s ObservableInt) Filter(f func(int) bool) ObservableInt {
	predicate := func(v interface{}) bool {
		return f(v.(int))
	}
	return s.adaptFilter(filters.Where(predicate))
}

// Last returns just the first element of the stream.
func (s ObservableInt) First() ObservableInt {
	return s.adaptFilter(filters.First())
}

// Last returns just the last element of the stream.
func (s ObservableInt) Last() ObservableInt {
	return s.adaptFilter(filters.Last())
}

// SkipLast skips the first N elements of the stream.
func (s ObservableInt) Skip(n int) ObservableInt {
	return s.adaptFilter(filters.Skip(n))
}

// SkipLast skips the last N elements of the stream.
func (s ObservableInt) SkipLast(n int) ObservableInt {
	return s.adaptFilter(filters.SkipLast(n))
}

// Take returns just the first N elements of the stream.
func (s ObservableInt) Take(n int) ObservableInt {
	return s.adaptFilter(filters.Take(n))
}

// TakeLast returns just the last N elements of the stream.
func (s ObservableInt) TakeLast(n int) ObservableInt {
	return s.adaptFilter(filters.TakeLast(n))
}

// IgnoreElements ignores elements of the stream and emits only the completion events.
func (s ObservableInt) IgnoreElements() ObservableInt {
	return s.adaptFilter(filters.IgnoreElements())
}

// IgnoreCompletion ignores the completion event of the stream and therefore returns a stream that never completes.
func (s ObservableInt) IgnoreCompletion() ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
		operator := func(next int, err error, completed bool) {
			if !completed {
				observer(next, err, completed)
			}
		}
		return s(operator)
	}
	return subscribe
}

func (s ObservableInt) One() ObservableInt {
	return s.adaptFilter(filters.One())
}

func (s ObservableInt) Replay(size int, duration time.Duration) ObservableInt {
	if size == 0 {
		size = MaxReplaySize
	}
	return s.adaptFilter(filters.Replay(size, duration))
}

func (s ObservableInt) Sample(duration time.Duration) ObservableInt {
	return s.adaptFilter(filters.Sample(duration))
}

// Debounce reduces subsequent duplicates to single items during a certain duration
func (s ObservableInt) Debounce(duration time.Duration) ObservableInt {
	return s.adaptFilter(filters.Debounce(duration))
}

/////////////////////////////////////////////////////////////////////////////
// COUNT
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) Count() ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
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
// PASSTHROUGH
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) Passthrough_0() ObservableInt {
	return s
}

func (s ObservableInt) Passthrough_1() ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
		return s(observer)
	}
	return subscribe
}

func (s ObservableInt) Passthrough_2() ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
		operator := func(next int, err error, complete bool) {
			observer(next, err, complete)
		}
		return s(operator)
	}
	return subscribe
}

func (s ObservableInt) Passthrough_3() ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
		unsubscribers := &unsubscriber.Collection{}
		operator := func(next int, err error, complete bool) {
			observer(next, err, complete)
		}
		subscribeOn := func(scheduler Scheduler) {
			// if SubscribeOn is in parent chain, the operator will
			// already be running once we get the subscribable back.
			subscribable := s(operator)
			// Set subscribable in the unsubscribers so Unsubscribe() will
			// propagate correctly from the subscriber to the observable.
			if !unsubscribers.Set(subscribable) {
				return
			}
			// Execute the subscribable observable on a scheduler
			subscribable.SubscribeOn(scheduler)
		}
		return &struct {
			SubscribeOnFunc
			Unsubscriber
		}{subscribeOn, unsubscribers}
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// CONCAT
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) Concat(other ...ObservableInt) ObservableInt {
	if len(other) == 0 {
		return s
	}
	subscribe := func(observer IntObserverFunc) Subscribable {
		unsubscribers := &unsubscriber.Collection{}

		var index int
		var maxIndex = 1 + len(other) - 1

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

		subscribeOn := func(scheduler Scheduler) {
			// Execute the first observable on the passed in scheduler.
			wg.Add(1)
			subscribable := s(operator)
			if !unsubscribers.Set(subscribable) {
				return
			}
			subscribable.SubscribeOn(scheduler)
			wg.Wait()

			for i, observable := range other {
				index = i + 1
				wg.Add(1)
				unsubscriber := observable.Subscribe(operator)
				if !unsubscribers.Set(unsubscriber) {
					return
				}
				wg.Wait()
				if index == maxIndex {
					return
				}
			}
		}

		return &struct {
			SubscribeOnFunc
			Unsubscriber
		}{subscribeOn, unsubscribers}
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// MERGE
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) merge(other []ObservableInt, delayError bool) ObservableInt {
	if len(other) == 0 {
		return s
	}
	subscribe := func(observer IntObserverFunc) Subscribable {
		unsubscribers := &unsubscriber.Collection{}

		var (
			olock     sync.Mutex
			merged    int
			count     = 1 + len(other)
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

		subscribeOn := func(scheduler Scheduler) {
			subscribable := s(operator)
			if !unsubscribers.Add(subscribable) {
				return
			}
			subscribable.SubscribeOn(scheduler)
			for _, observable := range other {
				unsubscriber := observable.Subscribe(operator)
				if !unsubscribers.Add(unsubscriber) {
					return
				}
			}
		}

		return &struct {
			SubscribeOnFunc
			Unsubscriber
		}{subscribeOn, unsubscribers}
	}
	return subscribe
}

// Merge an arbitrary number of observables with this one.
// An error from any of the observables will terminate the merged stream.
func (s ObservableInt) Merge(other ...ObservableInt) ObservableInt {
	return s.merge(other, false)
}

// Merge an arbitrary number of observables with this one.
// Any error will be deferred until all observables terminate.
func (s ObservableInt) MergeDelayError(other ...ObservableInt) ObservableInt {
	return s.merge(other, true)
}

/////////////////////////////////////////////////////////////////////////////
// CATCH
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) Catch(catch ObservableInt) ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
		unsubscribers := &unsubscriber.Collection{}

		throwChan := make(chan error, 1)
		operator := func(next int, err error, completed bool) {
			switch {
			case err != nil:
				throwChan <- err
			case completed:
				observer.Complete()
				throwChan <- nil
			default:
				observer.Next(next)
			}
		}

		catcher := func() {
			if err := <-throwChan; err != nil {
				unsubscriber := catch.Subscribe(observer)
				if !unsubscribers.Set(unsubscriber) {
					return
				}
			}
		}

		subscribeOn := func(scheduler Scheduler) {
			go catcher()

			subscribable := s(operator)
			if !unsubscribers.Set(subscribable) {
				return
			}
			subscribable.SubscribeOn(scheduler)
		}

		unsubscribe := func() {
			throwChan <- nil
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return &struct {
			SubscribeOnFunc
			Unsubscriber
		}{subscribeOn, unsubscribers}
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// RETRY
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) Retry() ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
		unsubscribers := &unsubscriber.Collection{}

		throwChan := make(chan error, 1)
		operator := func(next int, err error, completed bool) {
			switch {
			case err != nil:
				throwChan <- err
			case completed:
				observer.Complete()
				throwChan <- nil
			default:
				observer.Next(next)
			}
		}

		catcher := func(scheduler Scheduler) {
			for err := range throwChan {
				if err == nil {
					return
				}
				subscribable := s(operator)
				if !unsubscribers.Set(subscribable) {
					return
				}
				subscribable.SubscribeOn(scheduler)
			}
		}

		subscribeOn := func(scheduler Scheduler) {
			go catcher(scheduler)

			subscribable := s(operator)
			if !unsubscribers.Set(subscribable) {
				return
			}
			subscribable.SubscribeOn(scheduler)
		}

		unsubscribe := func() {
			throwChan <- nil
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return &struct {
			SubscribeOnFunc
			Unsubscriber
		}{subscribeOn, unsubscribers}
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// DO
/////////////////////////////////////////////////////////////////////////////

// Do applies a function for each value passing through the stream.
func (s ObservableInt) Do(f func(next int)) ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
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
func (s ObservableInt) DoOnError(f func(err error)) ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
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
func (s ObservableInt) DoOnComplete(f func()) ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
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
func (s ObservableInt) Finally(f func()) ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
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

func (s ObservableInt) Reduce(initial int, reducer func(int, int) int) ObservableInt {
	value := initial
	subscribe := func(observer IntObserverFunc) Subscribable {
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

func (s ObservableInt) Scan(initial int, f func(int, int) int) ObservableInt {
	value := initial
	subscribe := func(observer IntObserverFunc) Subscribable {
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

func (s ObservableInt) Timeout(timeout time.Duration) ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
		unsubscribers := &unsubscriber.Collection{}

		unsubchan := make(chan struct{})
		errchan := make(chan error)
		nextchan := make(chan int)
		deliver := func() {
			for {
				select {
				case <-unsubchan:
					return
				case <-time.After(timeout):
					observer(zeroInt, ErrTimeout, false)
					unsubscribers.Unsubscribe()
					return
				case err := <-errchan:
					observer(zeroInt, err, err == nil)
					unsubscribers.Unsubscribe()
					return
				case next := <-nextchan:
					observer(next, nil, false)
				}
			}
		}

		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				errchan <- err
			} else {
				nextchan <- next
			}
		}

		subscribeOn := func(scheduler Scheduler) {
			go deliver()

			subscribable := s(operator)
			if !unsubscribers.Set(subscribable) {
				return
			}
			subscribable.SubscribeOn(scheduler)
		}

		unsubscribe := func() {
			close(unsubchan)
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return &struct {
			SubscribeOnFunc
			Unsubscriber
		}{subscribeOn, unsubscribers}
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// FORK
/////////////////////////////////////////////////////////////////////////////

// Fork replicates each event from the parent to every subscriber of the fork.
// This allows multiple subscriptions to a single observable.
func (s ObservableInt) Fork() ObservableInt {
	var observers struct {
		sync.Mutex
		items []IntObserverFunc
	}

	addObserver := func(observer IntObserverFunc) int {
		observers.Lock()
		defer observers.Unlock()
		index := len(observers.items)
		observers.items = append(observers.items, observer)
		return index
	}

	removeObserver := func(index int) {
		observers.Lock()
		defer observers.Unlock()
		observers.items[index] = nil
	}

	eachObserver := func(f func(observer IntObserverFunc)) {
		observers.Lock()
		defer observers.Unlock()
		for _, observer := range observers.items {
			if observer != nil {
				f(observer)
			}
		}
	}

	operator := func(next int, err error, completed bool) {
		eachObserver(func(observer IntObserverFunc) {
			observer(next, err, completed)
		})
	}

	subscribe := func(observer IntObserverFunc) Subscribable {
		unsubscribers := &unsubscriber.Collection{}
		index := addObserver(observer)

		subscribeOn := func(scheduler Scheduler) {
		}

		unsubscribe := func() {
			removeObserver(index)
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return &struct {
			SubscribeOnFunc
			Unsubscriber
		}{subscribeOn, unsubscribers}
	}

	// Subscribes immediately on creation of the observable.
	// So before anybody actually subscribed.
	s.Subscribe(operator)

	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// PUBLISH
/////////////////////////////////////////////////////////////////////////////

type ConnectableInt struct {
	ObservableInt
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
func (s ObservableInt) Publish() ConnectableInt {
	var observers struct {
		sync.Mutex
		items []IntObserverFunc
	}

	addObserver := func(observer IntObserverFunc) int {
		observers.Lock()
		defer observers.Unlock()
		index := len(observers.items)
		observers.items = append(observers.items, observer)
		return index
	}

	removeObserver := func(index int) {
		observers.Lock()
		defer observers.Unlock()
		observers.items[index] = nil
	}

	eachObserver := func(f func(observer IntObserverFunc)) {
		observers.Lock()
		defer observers.Unlock()
		for _, observer := range observers.items {
			if observer != nil {
				f(observer)
			}
		}
	}

	operator := func(next int, err error, completed bool) {
		eachObserver(func(observer IntObserverFunc) {
			observer(next, err, completed)
		})
	}

	subscribe := func(observer IntObserverFunc) Subscribable {
		unsubscribers := &unsubscriber.Collection{}
		index := addObserver(observer)

		subscribeOn := func(scheduler Scheduler) {
		}

		unsubscribe := func() {
			removeObserver(index)
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return &struct {
			SubscribeOnFunc
			Unsubscriber
		}{subscribeOn, unsubscribers}
	}

	// Connect our operator observer function to the source start receiving
	// values from the source and forward them to our subscribers.
	connect := func() Unsubscriber {
		return s.Subscribe(operator)
	}

	return ConnectableInt{ObservableInt: subscribe, connect: connect}
}

/////////////////////////////////////////////////////////////////////////////
// MATHEMATICAL
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) Average() ObservableInt {
	var sum int
	var count int
	subscribe := func(observer IntObserverFunc) Subscribable {
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

func (s ObservableInt) Sum() ObservableInt {
	var sum int
	subscribe := func(observer IntObserverFunc) Subscribable {
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

func (s ObservableInt) Min() ObservableInt {
	started := false
	var min int
	subscribe := func(observer IntObserverFunc) Subscribable {
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

func (s ObservableInt) Max() ObservableInt {
	started := false
	var max int
	subscribe := func(observer IntObserverFunc) Subscribable {
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
// MAP (int,int)
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) Map(f func(int) int) ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
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

/////////////////////////////////////////////////////////////////////////////
// FLATMAP (int,int)
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) FlatMap(f func(int) ObservableInt) ObservableInt {
	subscribe := func(observer IntObserverFunc) Subscribable {
		unsubscribers := &unsubscriber.Collection{}

		var wait struct {
			sync.Mutex
			Condition *sync.Cond
			Accum     int
		}
		wait.Condition = sync.NewCond(&wait.Mutex)

		waitForZero := func() bool {
			wait.Lock()
			defer wait.Unlock()
			for wait.Accum > 0 {
				wait.Condition.Wait()
			}
			return (wait.Accum == 0)
		}

		waitAdd := func() {
			wait.Lock()
			defer wait.Unlock()
			wait.Accum += 1
			wait.Condition.Broadcast()
		}

		waitDone := func() {
			wait.Lock()
			defer wait.Unlock()
			wait.Accum -= 1
			wait.Condition.Broadcast()
		}

		waitCancel := func() {
			wait.Lock()
			defer wait.Unlock()
			wait.Accum = -1
			wait.Condition.Broadcast()
		}

		var lock sync.Mutex
		flatten := func(next int, err error, completed bool) {
			lock.Lock()
			defer lock.Unlock()
			// Finally
			if err != nil || completed {
				waitDone()
			}
			// IgnoreCompletion
			if !completed {
				observer(next, err, completed)
			}
		}

		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				if waitForZero() {
					observer(zeroInt, err, err == nil)
				}
			} else {
				waitAdd()
				unsubscriber := f(next).Subscribe(flatten)
				if !unsubscribers.Add(unsubscriber) {
					return
				}
			}
		}

		subscribeOn := func(scheduler Scheduler) {
			subscribable := s(operator)
			if !unsubscribers.Add(subscribable) {
				return
			}
			subscribable.SubscribeOn(scheduler)
		}

		unsubscribe := func() {
			waitCancel()
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return &struct {
			SubscribeOnFunc
			Unsubscriber
		}{subscribeOn, unsubscribers}
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// MAP (int,float64)
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) MapFloat64(f func(int) float64) ObservableFloat64 {
	subscribe := func(observer Float64ObserverFunc) Subscribable {
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

/////////////////////////////////////////////////////////////////////////////
// FLATMAP (int,float64)
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) FlatMapFloat64(f func(int) ObservableFloat64) ObservableFloat64 {
	subscribe := func(observer Float64ObserverFunc) Subscribable {
		unsubscribers := &unsubscriber.Collection{}

		var wait struct {
			sync.Mutex
			Condition *sync.Cond
			Accum     int
		}
		wait.Condition = sync.NewCond(&wait.Mutex)

		waitForZero := func() bool {
			wait.Lock()
			defer wait.Unlock()
			for wait.Accum > 0 {
				wait.Condition.Wait()
			}
			return (wait.Accum == 0)
		}

		waitAdd := func() {
			wait.Lock()
			defer wait.Unlock()
			wait.Accum++
			wait.Condition.Broadcast()
		}

		waitDone := func() {
			wait.Lock()
			defer wait.Unlock()
			wait.Accum--
			wait.Condition.Broadcast()
		}

		waitCancel := func() {
			wait.Lock()
			defer wait.Unlock()
			wait.Accum = -1
			wait.Condition.Broadcast()
		}

		var lock sync.Mutex
		flatten := func(next float64, err error, completed bool) {
			lock.Lock()
			defer lock.Unlock()
			// Finally
			if err != nil || completed {
				waitDone()
			}
			// IgnoreCompletion
			if !completed {
				observer(next, err, completed)
			}
		}

		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				if waitForZero() {
					observer(zeroFloat64, err, err == nil)
				}
			} else {
				waitAdd()
				unsubscriber := f(next).Subscribe(flatten)
				if !unsubscribers.Add(unsubscriber) {
					return
				}

			}
		}

		subscribeOn := func(scheduler Scheduler) {
			subscribable := s(operator)
			if !unsubscribers.Add(subscribable) {
				return
			}
			subscribable.SubscribeOn(scheduler)
		}

		unsubscribe := func() {
			waitCancel()
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return &struct {
			SubscribeOnFunc
			Unsubscriber
		}{subscribeOn, unsubscribers}
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// MAP (int,string)
/////////////////////////////////////////////////////////////////////////////

func (s ObservableInt) MapString(f func(int) string) ObservableString {
	subscribe := func(observer StringObserverFunc) Subscribable {
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

////////////////////////////////////////////////////////
// ObservableFloat64
////////////////////////////////////////////////////////

type ObservableFloat64 func(Float64ObserverFunc) Subscribable

type Float64ObserverFunc func(float64, error, bool)

var zeroFloat64 float64

func (s ObservableFloat64) Subscribe(observer Float64ObserverFunc) Unsubscriber {
	subscribable := s(observer)
	subscribable.SubscribeOn(schedulers.GoroutineScheduler)
	return &struct{ Unsubscriber }{subscribable}
}

////////////////////////////////////////////////////////
// ObservableString
////////////////////////////////////////////////////////

type ObservableString func(StringObserverFunc) Subscribable

type StringObserverFunc func(string, error, bool)

var zeroString string

func (s ObservableString) Subscribe(observer StringObserverFunc) Unsubscriber {
	subscribable := s(observer)
	subscribable.SubscribeOn(schedulers.GoroutineScheduler)
	return &struct{ Unsubscriber }{subscribable}
}
