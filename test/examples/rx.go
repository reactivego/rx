package examples

import (
	"errors"
	"rxgo/filters"
	"rxgo/scheduler"
	"rxgo/unsubscriber"
	"sync"
	"time"
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

/////////////////////////////////////////////////////////////////////////////
// FROM
/////////////////////////////////////////////////////////////////////////////

// CreateInt calls f(observer) to produce values for a stream of ints.
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

func FromIntChannelWithError(data <-chan int, errs <-chan error) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		for {
			select {
			case next, ok := <-data:
				if observer.Unsubscribed() {
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
				if observer.Unsubscribed() {
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

func Interval(interval time.Duration) ObservableInt {
	return CreateInt(func(observer IntObserver) {
		for i := 0; ; i++ {
			time.Sleep(interval)
			if observer.Unsubscribed() {
				return
			}
			observer.Next(i)
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
	observable := func(observer IntObserverFunc, _ Scheduler, unsubscriber Unsubscriber) {
		// override the scheduler
		o(observer, scheduler, unsubscriber)
	}
	return observable
}

// Subscribe an observer function to the observable.
// This method returns an Unsubscriber. You can always call Unsubscribe
// on the returned Unsubscriber to make sure all resources (goroutines)
// are deallocated. If the subscription terminates naturally because
// of an error or when the stream of data is complete then the Unsubscribe
// method on the Unsubscriber is automatically called internally.
func (o ObservableInt) Subscribe(observer IntObserverFunc) Unsubscriber {
	unsubscriber := unsubscriber.New()
	operator := func(next int, err error, completed bool) {
		observer(next, err, completed)
		if err != nil || completed {
			unsubscriber.Unsubscribe()
		}
	}
	o(operator, scheduler.Goroutines, unsubscriber)
	return unsubscriber
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
	return o.adaptFilter(filters.IgnoreCompletion())
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
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		count := 0
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(count, nil, false)
				observer(zeroInt, err, completed)
			}
			count++
		}
		o(operator, scheduler, unsubscriber)
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
	observable := func(observer IntObserverFunc, sched Scheduler, unsub Unsubscriber) {
		var index int
		var maxIndex = 1 + len(other) - 1

		operator := func(next int, err error, completed bool) {
			switch {
			case err != nil:
				index = maxIndex
				observer.Error(err)
			case completed:
				if index == maxIndex {
					observer.Complete()
				}
			default:
				observer.Next(next)
			}
		}

		// Execute the concat on the passed in scheduler.
		sched.Schedule(func() {
			o(operator, scheduler.Immediate, unsub)
			if index == maxIndex {
				return
			}
			for i, observable := range other {
				index = i + 1
				observable(operator, scheduler.Immediate, unsub)
				if index == maxIndex {
					return
				}
			}
		})
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
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
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

		o(operator, scheduler, unsubscriber)
		for _, o := range other {
			o(operator, scheduler, unsubscriber)
		}
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
	observable := func(observer IntObserverFunc, sched Scheduler, unsub Unsubscriber) {
		operator := func(next int, err error, completed bool) {
			if err != nil {
				catch(observer, scheduler.Immediate, unsub)
			} else {
				observer(next, err, completed)
			}
		}
		o(operator, sched, unsub)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// RETRY
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Retry() ObservableInt {
	observable := func(observer IntObserverFunc, sched Scheduler, unsub Unsubscriber) {
		sched.Schedule(func() {
			if !unsub.Unsubscribed() {
				var done bool
				operator := func(next int, err error, completed bool) {
					if err == nil {
						observer(next, err, completed)
						done = completed
					}
				}
				for !unsub.Unsubscribed() && !done {
					o(operator, scheduler.Immediate, unsub)
				}
			}
		})
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// DO
/////////////////////////////////////////////////////////////////////////////

// Do applies a function for each value passing through the stream.
func (o ObservableInt) Do(f func(next int)) ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		operator := func(next int, err error, completed bool) {
			if err == nil && !completed {
				f(next)
			}
			observer(next, err, completed)
		}
		o(operator, scheduler, unsubscriber)
	}
	return observable
}

// DoOnError applies a function for any error on the stream.
func (o ObservableInt) DoOnError(f func(err error)) ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		operator := func(next int, err error, completed bool) {
			if err != nil {
				f(err)
			}
			observer(next, err, completed)
		}
		o(operator, scheduler, unsubscriber)
	}
	return observable
}

// DoOnComplete applies a function when the stream completes.
func (o ObservableInt) DoOnComplete(f func()) ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		operator := func(next int, err error, completed bool) {
			if completed {
				f()
			}
			observer(next, err, completed)
		}
		o(operator, scheduler, unsubscriber)
	}
	return observable
}

// Finally applies a function for any error or completion on the stream.
// This doesn't expose whether this was an error or a completion.
func (o ObservableInt) Finally(f func()) ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				f()
			}
			observer(next, err, completed)
		}
		o(operator, scheduler, unsubscriber)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// REDUCE
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Reduce(initial int, reducer func(int, int) int) ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		value := initial
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(value, nil, false)
				observer(zeroInt, err, completed)
			} else {
				value = reducer(value, next)
			}
		}
		o(operator, scheduler, unsubscriber)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// SCAN
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Scan(initial int, f func(int, int) int) ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		value := initial
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(zeroInt, err, completed)
			} else {
				value = f(value, next)
				observer(value, nil, false)
			}
		}
		o(operator, scheduler, unsubscriber)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// TIMEOUT
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Timeout(timeout time.Duration) ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {

		if scheduler.Asynchronous() {
			// Asynchronous, so all scheduled tasks are executed in parallel.
			unsubscriber := unsubscriber.AddChild()

			deadline := time.NewTimer(timeout)

			var lock sync.Mutex
			operator := func(next int, err error, completed bool) {
				lock.Lock()
				defer lock.Unlock()

				if deadline.Stop() {
					// next data received before deadline expired.
					observer(next, err, completed)
					deadline.Reset(timeout)
				} else {
					// unfortunately the deadline has expired.
					observer.Error(ErrTimeout)
					unsubscriber.Unsubscribe()
				}
			}

			scheduler.Schedule(func() {
				<-deadline.C
				IntObserverFunc(operator).Error(ErrTimeout)
			})

			o(operator, scheduler, unsubscriber)

		} else {
			// Synchronous, meaning all Scheduled tasks are executed in sequence.
			unsubscriber := unsubscriber.AddChild()
			lastTime := time.Now()
			operator := func(next int, err error, completed bool) {
				if time.Since(lastTime) > timeout {
					unsubscriber.Unsubscribe()
					observer.Error(ErrTimeout)
				} else {
					observer(next, err, completed)
					lastTime = time.Now()
				}
			}
			o(operator, scheduler, unsubscriber)
		}
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// FORK
/////////////////////////////////////////////////////////////////////////////

// Fork replicates each event from the parent to every observer of the fork.
// This allows multiple subscriptions to a single observable.
// This only works for hot observables, cold observables may drain out to the
// first connected observable even before a subsequent observer subscribes.
func (o ObservableInt) Fork() ObservableInt {
	fork := o.Publish()
	// TODO: Nobody can actually Unsubscribe from this thing....
	fork.Connect()
	return fork.ObservableInt
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

	type channel struct {
		data chan int
		errs chan error
	}

	var channels struct {
		sync.Mutex
		items []*channel
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

	operator := func(next int, err error, completed bool) {
		switch {
		case err != nil:
			broadcastError(err)
		case completed:
			removeAllChannels()
		default:
			broadcastData(next)
		}
	}

	// Connect our operator observer function to the source start receiving
	// values from the source and forward them to our subscribers.
	connect := func() Unsubscriber {
		// TODO: Can Connect() Unsubscribe() Connect() sequence be used to pause/restart the stream?
		return o.Subscribe(operator)
	}

	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		data := make(chan int, 1)
		errs := make(chan error, 1)
		index := addChannel(data, errs)
		unsubscriber.OnUnsubscribe(func() {
			removeChannel(index)
		})
		FromIntChannelWithError(data, errs)(observer, scheduler, unsubscriber)
	}

	return ConnectableInt{ObservableInt: observable, connect: connect}
}

/////////////////////////////////////////////////////////////////////////////
// MATHEMATICAL
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) Average() ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		var sum int
		var count int
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(sum/count, nil, false)
				observer(zeroInt, err, completed)
			} else {
				sum += next
				count++
			}
		}
		o(operator, scheduler, unsubscriber)
	}
	return observable
}

func (o ObservableInt) Sum() ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		var sum int
		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				observer(sum, nil, false)
				observer(zeroInt, err, completed)
			} else {
				sum += next
			}
		}
		o(operator, scheduler, unsubscriber)
	}
	return observable
}

func (o ObservableInt) Min() ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		started := false
		var min int
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
		o(operator, scheduler, unsubscriber)
	}
	return observable
}

func (o ObservableInt) Max() ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		started := false
		var max int
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
		o(operator, scheduler, unsubscriber)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// MAP (int,int)
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) MapInt(f func(int) int) ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		operator := func(next int, err error, completed bool) {
			var mapped int
			if err == nil && !completed {
				mapped = f(next)
			}
			observer(mapped, err, completed)
		}
		o(operator, scheduler, unsubscriber)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// MAP (int,float64)
/////////////////////////////////////////////////////////////////////////////

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

/////////////////////////////////////////////////////////////////////////////
// MAP (int,string)
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) MapString(f func(int) string) ObservableString {
	observable := func(observer StringObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		operator := func(next int, err error, completed bool) {
			var mapped string
			if err == nil && !completed {
				mapped = f(next)
			}
			observer(mapped, err, completed)
		}
		o(operator, scheduler, unsubscriber)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// FLATMAP (int,int)
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) FlatMapInt(f func(int) ObservableInt) ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		var wait struct {
			sync.Mutex
			Condition *sync.Cond
			Accum     int

			Add         func()
			Done        func()
			Cancel      func()
			WaitForZero func() bool
		}
		wait.Condition = sync.NewCond(&wait.Mutex)

		wait.Add = func() {
			wait.Lock()
			defer wait.Unlock()
			wait.Accum++
			wait.Condition.Broadcast()
		}

		wait.Done = func() {
			wait.Lock()
			defer wait.Unlock()
			wait.Accum--
			wait.Condition.Broadcast()
		}

		wait.Cancel = func() {
			wait.Lock()
			defer wait.Unlock()
			wait.Accum = -1
			wait.Condition.Broadcast()
		}

		wait.WaitForZero = func() bool {
			wait.Lock()
			defer wait.Unlock()
			for wait.Accum > 0 {
				wait.Condition.Wait()
			}
			return (wait.Accum == 0)
		}

		var lock sync.Mutex
		flatten := func(next int, err error, completed bool) {
			lock.Lock()
			defer lock.Unlock()
			// Finally
			if err != nil || completed {
				wait.Done()
			}
			// IgnoreCompletion
			if !completed {
				observer(next, err, completed)
			}
		}

		unsubscriber.OnUnsubscribe(wait.Cancel)

		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				if wait.WaitForZero() {
					observer(zeroInt, err, err == nil)
				}
			} else {
				wait.Add()
				f(next)(flatten, scheduler, unsubscriber)
			}
		}

		o(operator, scheduler, unsubscriber)
	}
	return observable
}

/////////////////////////////////////////////////////////////////////////////
// FLATMAP (int,float64)
/////////////////////////////////////////////////////////////////////////////

func (o ObservableInt) FlatMapFloat64(f func(int) ObservableFloat64) ObservableFloat64 {
	observable := func(observer Float64ObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		var wait struct {
			sync.Mutex
			Condition *sync.Cond
			Accum     int

			Add         func()
			Done        func()
			Cancel      func()
			WaitForZero func() bool
		}
		wait.Condition = sync.NewCond(&wait.Mutex)

		wait.Add = func() {
			wait.Lock()
			defer wait.Unlock()
			wait.Accum++
			wait.Condition.Broadcast()
		}

		wait.Done = func() {
			wait.Lock()
			defer wait.Unlock()
			wait.Accum--
			wait.Condition.Broadcast()
		}

		wait.Cancel = func() {
			wait.Lock()
			defer wait.Unlock()
			wait.Accum = -1
			wait.Condition.Broadcast()
		}

		wait.WaitForZero = func() bool {
			wait.Lock()
			defer wait.Unlock()
			for wait.Accum > 0 {
				wait.Condition.Wait()
			}
			return (wait.Accum == 0)
		}

		var lock sync.Mutex
		flatten := func(next float64, err error, completed bool) {
			lock.Lock()
			defer lock.Unlock()
			// Finally
			if err != nil || completed {
				wait.Done()
			}
			// IgnoreCompletion
			if !completed {
				observer(next, err, completed)
			}
		}

		unsubscriber.OnUnsubscribe(wait.Cancel)

		operator := func(next int, err error, completed bool) {
			if err != nil || completed {
				if wait.WaitForZero() {
					observer(zeroFloat64, err, err == nil)
				}
			} else {
				wait.Add()
				f(next)(flatten, scheduler, unsubscriber)
			}
		}

		o(operator, scheduler, unsubscriber)
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
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		o(observer, scheduler, unsubscriber)
	}
	return observable
}

func (o ObservableInt) Passthrough_2() ObservableInt {
	observable := func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		operator := func(next int, err error, complete bool) {
			observer(next, err, complete)
		}
		o(operator, scheduler, unsubscriber)
	}
	return observable
}

func (o ObservableInt) Passthrough_3() ObservableInt {
	return func(observer IntObserverFunc, scheduler Scheduler, unsubscriber Unsubscriber) {
		o(func(next int, err error, complete bool) {
			observer(next, err, complete)
		}, scheduler, unsubscriber)
	}
}

////////////////////////////////////////////////////////
// ObservableFloat64
////////////////////////////////////////////////////////

type Float64ObserverFunc func(float64, error, bool)

var zeroFloat64 float64

type ObservableFloat64 func(Float64ObserverFunc, Scheduler, Unsubscriber)

func (o ObservableFloat64) Subscribe(observer Float64ObserverFunc) Unsubscriber {
	unsubscriber := unsubscriber.New()
	operator := func(next float64, err error, completed bool) {
		observer(next, err, completed)
		if err != nil || completed {
			unsubscriber.Unsubscribe()
		}
	}
	o(operator, scheduler.Goroutines, unsubscriber)
	return unsubscriber
}

////////////////////////////////////////////////////////
// ObservableString
////////////////////////////////////////////////////////

type StringObserverFunc func(string, error, bool)

var zeroString string

type ObservableString func(StringObserverFunc, Scheduler, Unsubscriber)

func (o ObservableString) Subscribe(observer StringObserverFunc) Unsubscriber {
	unsubscriber := unsubscriber.New()
	operator := func(next string, err error, completed bool) {
		observer(next, err, completed)
		if err != nil || completed {
			unsubscriber.Unsubscribe()
		}
	}
	o(operator, scheduler.Goroutines, unsubscriber)
	return unsubscriber
}
