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

// Subscribable is a struct returned as a result of calling
// the observable function. See e.g. type ObservableInt.
type Subscribable struct {
	SubscribeOn func(scheduler Scheduler)
	Unsubscriber
}

////////////////////////////////////////////////////////
// IntObserver
////////////////////////////////////////////////////////

type IntObserver interface {
	Next(int)
	Error(error)
	Complete()
	Unsubscriber
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

////////////////////////////////////////////////////////
// ObservableInt
////////////////////////////////////////////////////////

// Every observable is essentially a function taking
// an Observer (function) and returning a Subscribable.
type ObservableInt func(IntObserverFunc) Subscribable

/////////////////////////////////////////////////////////////////////////////
// FROM
/////////////////////////////////////////////////////////////////////////////

// CreateInt calls f(observer) to produce values for a stream of ints.
func CreateInt(f func(IntObserver)) ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
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
					observer := &struct {
						IntObserverFunc
						Unsubscriber
					}{operation, unsubscriber}
					f(observer)
				}
			})
		}

		return Subscribable{subscribeOn, unsubscriber}
	}
	return observable
}

func EmptyInt() ObservableInt {
	return CreateInt(func(observer IntObserver) {
		observer.Complete()
	})
}

func NeverInt() ObservableInt {
	return CreateInt(func(observer IntObserver) {
	})
}

func ThrowInt(err error) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		observer.Error(err)
	})
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

func FromIntChannel(ch <-chan int) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		for next := range ch {
			if observer.Unsubscribed() {
				return
			}
			observer.Next(next)
		}
		observer.Complete()
	})
}

func Interval(interval time.Duration) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		i := 0
		for {
			time.Sleep(interval)
			if observer.Unsubscribed() {
				return
			}
			observer.Next(i)
			i++
		}
	})
}

func JustInt(element int) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		observer.Next(element)
		observer.Complete()
	})
}

func Range(start, count int) ObservableInt {
	end := start + count
	return CreateInt(func(observer IntObserver) {
		for i := start; i < end; i++ {
			if observer.Unsubscribed() {
				return
			}
			observer.Next(i)
		}
		observer.Complete()
	})
}

// Repeat value count times.
func RepeatInt(value int, count int) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		for i := 0; i < count; i++ {
			if observer.Unsubscribed() {
				return
			}
			observer.Next(value)
		}
		observer.Complete()
	})
}

// StartInt is designed to be used with functions that return a
// (int, error) tuple.
//
// If the error is non-nil the returned ObservableInt will be that error,
// otherwise it will be a single-value stream of int.
func StartInt(f func() (int, error)) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		if next, err := f(); err != nil {
			observer.Error(err)
		} else {
			observer.Next(next)
			observer.Complete()
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

func (o ObservableInt) SubscribeOn(scheduler Scheduler) ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
		subscribable := o(observer)
		subscribable.SubscribeOn(scheduler)
		subscribable.SubscribeOn = func(scheduler Scheduler) { /* ignore */ }
		return subscribable
	}
	return observable
}

func (o ObservableInt) Subscribe(observer IntObserverFunc) Unsubscriber {
	subscribable := o(observer)
	subscribable.SubscribeOn(schedulers.GoroutineScheduler)
	return subscribable.Unsubscriber
}

func (o ObservableInt) SubscribeNext(f func(v int)) Unsubscriber {
	operator := func(next int, err error, completed bool) {
		if err == nil && !completed {
			f(next)
		}
	}
	return o.Subscribe(operator)
}

// Wait for completion of the stream and return any error.
func (o ObservableInt) Wait() error {
	doneChan := make(chan error)
	operator := func(next int, err error, completed bool) {
		if err != nil || completed {
			doneChan <- err
		}
	}
	o.Subscribe(operator)
	return <-doneChan
}

/////////////////////////////////////////////////////////////////////////////
// TO
/////////////////////////////////////////////////////////////////////////////

// ToOneWithError blocks until the stream emits exactly one value. Otherwise, it errors.
func (o ObservableInt) ToOneWithError() (v int, e error) {
	v = zeroInt
	errch := make(chan error, 1)
	o.One().Subscribe(func(next int, err error, completed bool) {
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
func (o ObservableInt) ToOne() int {
	value, _ := o.ToOneWithError()
	return value
}

// ToArrayWithError collects all values from the stream into an array,
// returning it and any error.
func (o ObservableInt) ToArrayWithError() (a []int, e error) {
	a = []int{}
	errch := make(chan error, 1)
	o.Subscribe(func(next int, err error, completed bool) {
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
func (o ObservableInt) ToArray() []int {
	out, _ := o.ToArrayWithError()
	return out
}

// ToChannelWithError returns next and error channels corresponding to the stream elements and
// any error. When the next channel closes, it may be either because of an error or
// because the observable completed. The next channel may be closed without having emitted any
// values. The error channel always emits a value to indicate the observable has finished.
// When the error channel emits nil then the observable completed without errors, otherwise
// the error channel emits the error. When the observable has finished both channels will be
// closed.
func (o ObservableInt) ToChannelWithError() (<-chan int, <-chan error) {
	nextch := make(chan int, 1)
	errch := make(chan error, 1)
	o.Subscribe(func(next int, err error, completed bool) {
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

func (o ObservableInt) ToChannel() <-chan int {
	ch, _ := o.ToChannelWithError()
	return ch
}

/////////////////////////////////////////////////////////////////////////////
// FILTERS
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) adaptFilter(filter filters.Filter) ObservableInt {
	observable := func(sink IntObserverFunc) Subscribable {
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
		return o(intToGenericSource)
	}
	return observable
}

// Distinct removes duplicate elements in the stream.
func (o ObservableInt) Distinct() ObservableInt {
	return o.adaptFilter(filters.Distinct())
}

// ElementAt yields the Nth element of the stream.
func (o ObservableInt) ElementAt(n int) ObservableInt {
	return o.adaptFilter(filters.ElementAt(n))
}

// Filter elements in the stream on a function.
func (o ObservableInt) Filter(f func(int) bool) ObservableInt {
	predicate := func(v interface{}) bool {
		return f(v.(int))
	}
	return o.adaptFilter(filters.Where(predicate))
}

// Last returns just the first element of the stream.
func (o ObservableInt) First() ObservableInt {
	return o.adaptFilter(filters.First())
}

// Last returns just the last element of the stream.
func (o ObservableInt) Last() ObservableInt {
	return o.adaptFilter(filters.Last())
}

// SkipLast skips the first N elements of the stream.
func (o ObservableInt) Skip(n int) ObservableInt {
	return o.adaptFilter(filters.Skip(n))
}

// SkipLast skips the last N elements of the stream.
func (o ObservableInt) SkipLast(n int) ObservableInt {
	return o.adaptFilter(filters.SkipLast(n))
}

// Take returns just the first N elements of the stream.
func (o ObservableInt) Take(n int) ObservableInt {
	return o.adaptFilter(filters.Take(n))
}

// TakeLast returns just the last N elements of the stream.
func (o ObservableInt) TakeLast(n int) ObservableInt {
	return o.adaptFilter(filters.TakeLast(n))
}

// IgnoreElements ignores elements of the stream and emits only the completion events.
func (o ObservableInt) IgnoreElements() ObservableInt {
	return o.adaptFilter(filters.IgnoreElements())
}

// IgnoreCompletion ignores the completion event of the stream and therefore returns a stream that never completes.
func (o ObservableInt) IgnoreCompletion() ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
		operator := func(next int, err error, completed bool) {
			if !completed {
				observer(next, err, completed)
			}
		}
		return o(operator)
	}
	return observable
}

func (o ObservableInt) One() ObservableInt {
	return o.adaptFilter(filters.One())
}

func (o ObservableInt) Replay(size int, duration time.Duration) ObservableInt {
	if size == 0 {
		size = MaxReplaySize
	}
	return o.adaptFilter(filters.Replay(size, duration))
}

func (o ObservableInt) Sample(duration time.Duration) ObservableInt {
	return o.adaptFilter(filters.Sample(duration))
}

// Debounce reduces subsequent duplicates to single items during a certain duration
func (o ObservableInt) Debounce(duration time.Duration) ObservableInt {
	return o.adaptFilter(filters.Debounce(duration))
}

/////////////////////////////////////////////////////////////////////////////
// COUNT
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Count() ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
		count := 0
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(count, nil, false)
				observer(zeroInt, err, completed)
			}
			count++
		}
		return o(operator)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// PASSTHROUGH
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Passthrough_0() ObservableInt {
	return o
}

func (o ObservableInt) Passthrough_1() ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
		return o(observer)
	}
	return observable
}

func (o ObservableInt) Passthrough_2() ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
		operator := func(next int, err error, complete bool) {
			observer(next, err, complete)
		}
		return o(operator)
	}
	return observable
}

func (o ObservableInt) Passthrough_3() ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
		unsubscribers := &unsubscriber.Collection{}
		operator := func(next int, err error, complete bool) {
			observer(next, err, complete)
		}
		subscribeOn := func(scheduler Scheduler) {
			// if SubscribeOn is in parent chain, the operator will
			// already be running once we get the subscribable back.
			subscribable := o(operator)
			// Set subscribable in the unsubscribers so Unsubscribe() will
			// propagate correctly from the observer to the observable.
			if !unsubscribers.Set(subscribable) {
				return
			}
			// Execute the subscribable observable on a scheduler
			subscribable.SubscribeOn(scheduler)
		}
		return Subscribable{subscribeOn, unsubscribers}
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// CONCAT
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Concat(other ...ObservableInt) ObservableInt {
	if len(other) == 0 {
		return o
	}
	observable := func(observer IntObserverFunc) Subscribable {
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
			subscribable := o(operator)
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

		return Subscribable{subscribeOn, unsubscribers}
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// MERGE
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) merge(other []ObservableInt, delayError bool) ObservableInt {
	if len(other) == 0 {
		return o
	}
	observable := func(observer IntObserverFunc) Subscribable {
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
			subscribable := o(operator)
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

		return Subscribable{subscribeOn, unsubscribers}
	}
	return observable
}

// Merge an arbitrary number of observables with this one.
// An error from any of the observables will terminate the merged stream.
func (o ObservableInt) Merge(other ...ObservableInt) ObservableInt {
	return o.merge(other, false)
}

// Merge an arbitrary number of observables with this one.
// Any error will be deferred until all observables terminate.
func (o ObservableInt) MergeDelayError(other ...ObservableInt) ObservableInt {
	return o.merge(other, true)
}

/////////////////////////////////////////////////////////////////////////////
// CATCH
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Catch(catch ObservableInt) ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
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

			subscribable := o(operator)
			if !unsubscribers.Set(subscribable) {
				return
			}
			subscribable.SubscribeOn(scheduler)
		}

		unsubscribe := func() {
			throwChan <- nil
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return Subscribable{subscribeOn, unsubscribers}
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// RETRY
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Retry() ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
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
				subscribable := o(operator)
				if !unsubscribers.Set(subscribable) {
					return
				}
				subscribable.SubscribeOn(scheduler)
			}
		}

		subscribeOn := func(scheduler Scheduler) {
			go catcher(scheduler)

			subscribable := o(operator)
			if !unsubscribers.Set(subscribable) {
				return
			}
			subscribable.SubscribeOn(scheduler)
		}

		unsubscribe := func() {
			throwChan <- nil
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return Subscribable{subscribeOn, unsubscribers}
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// DO
/////////////////////////////////////////////////////////////////////////////

// Do applies a function for each value passing through the stream.
func (o ObservableInt) Do(f func(next int)) ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
		operator := func(next int, err error, completed bool) {
			if err == nil && !completed {
				f(next)
			}
			observer(next, err, completed)
		}
		return o(operator)
	}
	return observable
}

// DoOnError applies a function for any error on the stream.
func (o ObservableInt) DoOnError(f func(err error)) ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
		operator := func(next int, err error, completed bool) {
			if err != nil {
				f(err)
			}
			observer(next, err, completed)
		}
		return o(operator)
	}
	return observable
}

// DoOnComplete applies a function when the stream completes.
func (o ObservableInt) DoOnComplete(f func()) ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
		operator := func(next int, err error, completed bool) {
			if completed {
				f()
			}
			observer(next, err, completed)
		}
		return o(operator)
	}
	return observable
}

// Finally applies a function for any error or completion on the stream.
// This doesn't expose whether this was an error or a completion.
func (o ObservableInt) Finally(f func()) ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				f()
			}
			observer(next, err, completed)
		}
		return o(operator)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// REDUCE
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Reduce(initial int, reducer func(int, int) int) ObservableInt {
	value := initial
	observable := func(observer IntObserverFunc) Subscribable {
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(value, nil, false)
				observer(zeroInt, err, completed)
			} else {
				value = reducer(value, next)
			}
		}
		return o(operator)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// SCAN
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Scan(initial int, f func(int, int) int) ObservableInt {
	value := initial
	observable := func(observer IntObserverFunc) Subscribable {
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(zeroInt, err, completed)
			} else {
				value = f(value, next)
				observer(value, nil, false)
			}
		}
		return o(operator)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// TIMEOUT
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Timeout(timeout time.Duration) ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
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

			subscribable := o(operator)
			if !unsubscribers.Set(subscribable) {
				return
			}
			subscribable.SubscribeOn(scheduler)
		}

		unsubscribe := func() {
			close(unsubchan)
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return Subscribable{subscribeOn, unsubscribers}
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// FORK
/////////////////////////////////////////////////////////////////////////////

// Fork replicates each event from the parent to every observer of the fork.
// This allows multiple subscriptions to a single observable.
func (o ObservableInt) Fork() ObservableInt {
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

	observable := func(observer IntObserverFunc) Subscribable {
		unsubscribers := &unsubscriber.Collection{}
		index := addObserver(observer)

		subscribeOn := func(scheduler Scheduler) {
		}

		unsubscribe := func() {
			removeObserver(index)
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return Subscribable{subscribeOn, unsubscribers}
	}

	// Subscribes immediately on creation of the observable.
	// So before anybody actually subscribed.
	o.Subscribe(operator)

	return observable
}

/////////////////////////////////////////////////////////////////////////////
// PUBLISH
/////////////////////////////////////////////////////////////////////////////

type ConnectableInt struct {
	ObservableInt
	connect func() Unsubscriber
}

// Connect will observable to the parent observable to start receiving values.
// All values will then be passed on to the observers that subscribed to this
// connectable observable
func (c ConnectableInt) Connect() Unsubscriber {
	return c.connect()
}

// Publish creates a connectable observable that only starts emitting values
// after the Connect method is called on it.
func (o ObservableInt) Publish() ConnectableInt {
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

	observable := func(observer IntObserverFunc) Subscribable {
		unsubscribers := &unsubscriber.Collection{}
		index := addObserver(observer)

		subscribeOn := func(scheduler Scheduler) {
		}

		unsubscribe := func() {
			removeObserver(index)
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return Subscribable{subscribeOn, unsubscribers}
	}

	// Connect our operator observer function to the source start receiving
	// values from the source and forward them to our subscribers.
	connect := func() Unsubscriber {
		return o.Subscribe(operator)
	}

	return ConnectableInt{ObservableInt: observable, connect: connect}
}

/////////////////////////////////////////////////////////////////////////////
// MATHEMATICAL
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Average() ObservableInt {
	var sum int
	var count int
	observable := func(observer IntObserverFunc) Subscribable {
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(sum/count, nil, false)
				observer(zeroInt, err, completed)
			} else {
				sum += next
				count++
			}
		}
		return o(operator)
	}
	return observable
}

func (o ObservableInt) Sum() ObservableInt {
	var sum int
	observable := func(observer IntObserverFunc) Subscribable {
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(sum, nil, false)
				observer(zeroInt, err, completed)
			} else {
				sum += next
			}
		}
		return o(operator)
	}
	return observable
}

func (o ObservableInt) Min() ObservableInt {
	started := false
	var min int
	observable := func(observer IntObserverFunc) Subscribable {
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
		return o(operator)
	}
	return observable
}

func (o ObservableInt) Max() ObservableInt {
	started := false
	var max int
	observable := func(observer IntObserverFunc) Subscribable {
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
		return o(operator)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// MAP (int,int)
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Map(f func(int) int) ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
		operator := func(next int, err error, completed bool) {
			var mapped int
			if err == nil && !completed {
				mapped = f(next)
			}
			observer(mapped, err, completed)
		}
		return o(operator)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// FLATMAP (int,int)
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) FlatMap(f func(int) ObservableInt) ObservableInt {
	observable := func(observer IntObserverFunc) Subscribable {
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
			subscribable := o(operator)
			if !unsubscribers.Add(subscribable) {
				return
			}
			subscribable.SubscribeOn(scheduler)
		}

		unsubscribe := func() {
			waitCancel()
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return Subscribable{subscribeOn, unsubscribers}
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// MAP (int,float64)
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) MapFloat64(f func(int) float64) ObservableFloat64 {
	observable := func(observer Float64ObserverFunc) Subscribable {
		operator := func(next int, err error, completed bool) {
			var mapped float64
			if err == nil && !completed {
				mapped = f(next)
			}
			observer(mapped, err, completed)
		}
		return o(operator)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// FLATMAP (int,float64)
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) FlatMapFloat64(f func(int) ObservableFloat64) ObservableFloat64 {
	observable := func(observer Float64ObserverFunc) Subscribable {
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
			subscribable := o(operator)
			if !unsubscribers.Add(subscribable) {
				return
			}
			subscribable.SubscribeOn(scheduler)
		}

		unsubscribe := func() {
			waitCancel()
		}
		unsubscribers.OnUnsubscribe(unsubscribe)

		return Subscribable{subscribeOn, unsubscribers}
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// MAP (int,string)
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) MapString(f func(int) string) ObservableString {
	observable := func(observer StringObserverFunc) Subscribable {
		operator := func(next int, err error, completed bool) {
			var mapped string
			if err == nil && !completed {
				mapped = f(next)
			}
			observer(mapped, err, completed)
		}
		return o(operator)
	}
	return observable
}

////////////////////////////////////////////////////////
// ObservableFloat64
////////////////////////////////////////////////////////

type ObservableFloat64 func(Float64ObserverFunc) Subscribable

type Float64ObserverFunc func(float64, error, bool)

var zeroFloat64 float64

func (o ObservableFloat64) Subscribe(observer Float64ObserverFunc) Unsubscriber {
	subscribable := o(observer)
	subscribable.SubscribeOn(schedulers.GoroutineScheduler)
	return subscribable.Unsubscriber
}

////////////////////////////////////////////////////////
// ObservableString
////////////////////////////////////////////////////////

type ObservableString func(StringObserverFunc) Subscribable

type StringObserverFunc func(string, error, bool)

var zeroString string

func (o ObservableString) Subscribe(observer StringObserverFunc) Unsubscriber {
	subscribable := o(observer)
	subscribable.SubscribeOn(schedulers.GoroutineScheduler)
	return subscribable.Unsubscriber
}
