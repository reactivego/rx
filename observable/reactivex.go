/*
TODO


Changes to

Should Subscription interface be renamed to Unsubscriber and Dispose and Disposed
to Unsubscribe and Unsubscribed?

Should the CreateFunc be changed to return an OnUnsubscribed function that can be called by
a subscription implementation to let the observable know that the Subscriber is gone. Especially
when using go routines this could be used to close a channel so the go routine of the observable
terminates.


*/

package observable

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// ErrTimeout is delivered to an observer if the stream times out.
var ErrTimeout = errors.New("timeout")

// MaxReplaySize is the maximum size of a replay buffer. Can be modified.
var MaxReplaySize = 16384

////////////////////////////////////////////////////////
// Subscription
////////////////////////////////////////////////////////

type Subscription interface {
	Dispose()
	Disposed() bool
}

// SubscriptionEvents provides lifecycle event callbacks for a Subscription.
type SubscriptionEvents interface {
	OnUnsubscribe(func())
}

// Int32Subscription is an Int32 value with atomic operations implementing the Subscription interface

type Int32Subscription int32

func (t *Int32Subscription) Dispose() {
	atomic.StoreInt32((*int32)(t), 1)
}

func (t *Int32Subscription) Disposed() bool {
	return atomic.LoadInt32((*int32)(t)) == 1
}

// ClosedSubscription is a subscription that always reports that it is already closed.
type ClosedSubscription struct{}

func (ClosedSubscription) Dispose() {
}

func (ClosedSubscription) Disposed() bool {
	return true
}

// ChannelSubscription is implemented with a channel which is closed when unsubscribed.
type ChannelSubscription chan struct{}

func NewChannelSubscription() ChannelSubscription {
	return make(ChannelSubscription)
}

func (c ChannelSubscription) Dispose() {
	defer recover()
	close(c)
}

func (c ChannelSubscription) Disposed() bool {
	select {
	case _, ok := <-c:
		return !ok
	default:
		return false
	}
}

func (c ChannelSubscription) OnUnsubscribe(handler func()) {
	go func() {
		<-c
		handler()
	}()
}

// CallbackSubscription is implemented as a pointer to a callback function. Whenever Dispose()
// is called the callback is invoked. The method Disposed() returns true when the pointer to the
// callback function has been set to nil.
type CallbackSubscription func()

func (c *CallbackSubscription) Dispose() {
	if *c != nil {
		(*c)()
		*c = nil
	}
}

func (c *CallbackSubscription) Disposed() bool {
	return *c == nil
}

////////////////////////////////////////////////////////
// Observer
////////////////////////////////////////////////////////

type Observer func(interface{}, error, bool)

func (f Observer) Next(next interface{}) {
	f(next, nil, false)
}

func (f Observer) Error(err error) {
	f(nil, err, false)
}

func (f Observer) Complete() {
	f(nil, nil, true)
}

////////////////////////////////////////////////////////
//  Generic Filter Implementations
////////////////////////////////////////////////////////

type FilterFactory func(Observer) Observer

type FiltersNamespace struct{}

var filters FiltersNamespace

func (FiltersNamespace) Distinct() FilterFactory {
	factory := func(observer Observer) Observer {
		seen := map[interface{}]struct{}{}
		filter := func(next interface{}, err error, complete bool) {
			if err == nil && !complete {
				if _, ok := seen[next]; ok {
					return
				}
				seen[next] = struct{}{}
			}
			observer(next, err, complete)
		}
		return filter
	}
	return factory
}

func (FiltersNamespace) ElementAt(n int) FilterFactory {
	factory := func(observer Observer) Observer {
		i := 0
		filter := func(next interface{}, err error, complete bool) {
			switch {
			case err != nil:
				observer.Error(err)
			case complete:
				observer.Complete()
			default:
				if i == n {
					observer.Next(next)
				}
				i++
			}
		}
		return filter
	}
	return factory
}

func (FiltersNamespace) Filter(f func(next interface{}) bool) FilterFactory {
	factory := func(observer Observer) Observer {
		filter := func(next interface{}, err error, complete bool) {
			switch {
			case err != nil:
				observer.Error(err)
			case complete:
				observer.Complete()
			default:
				if f(next) {
					observer.Next(next)
				}
			}
		}
		return filter
	}
	return factory
}

func (FiltersNamespace) First() FilterFactory {
	factory := func(observer Observer) Observer {
		start := true
		filter := func(next interface{}, err error, complete bool) {
			switch {
			case err != nil:
				observer.Error(err)
			case complete:
				observer.Complete()
			default:
				if start {
					observer.Next(next)
				}
				start = false
			}
		}
		return filter
	}
	return factory
}

func (FiltersNamespace) Last() FilterFactory {
	factory := func(observer Observer) Observer {
		have := false
		var last interface{}
		filter := func(next interface{}, err error, complete bool) {
			switch {
			case err != nil:
				if have {
					observer.Next(last)
				}
				observer.Error(err)
			case complete:
				if have {
					observer.Next(last)
				}
				observer.Complete()
			default:
				last = next
				have = true
			}
		}
		return filter
	}
	return factory
}

func (FiltersNamespace) Skip(n int) FilterFactory {
	factory := func(observer Observer) Observer {
		i := 0
		filter := func(next interface{}, err error, complete bool) {
			if err != nil || complete || i >= n {
				observer(next, err, complete)
			}
			i++
		}
		return filter
	}
	return factory
}

func (FiltersNamespace) SkipLast(n int) FilterFactory {
	factory := func(observer Observer) Observer {
		read := 0
		write := 0
		n++
		buffer := make([]interface{}, n)
		filter := func(next interface{}, err error, complete bool) {
			switch {
			case err != nil:
				observer.Error(err)
			case complete:
				observer.Complete()
			default:
				buffer[write] = next
				write = (write + 1) % n
				if write == read {
					observer.Next(buffer[read])
					read = (read + 1) % n
				}
			}
		}
		return filter
	}
	return factory
}

func (FiltersNamespace) Take(n int) FilterFactory {
	factory := func(observer Observer) Observer {
		taken := 0
		filter := func(next interface{}, err error, complete bool) {
			if taken < n {
				observer(next, err, complete)
				if err == nil && !complete {
					taken++
					if taken >= n {
						observer.Complete()
					}
				}
			}
		}
		return filter
	}
	return factory
}

func (FiltersNamespace) TakeLast(n int) FilterFactory {
	factory := func(observer Observer) Observer {
		read := 0
		write := 0
		n++
		buffer := make([]interface{}, n)
		filter := func(next interface{}, err error, complete bool) {
			switch {
			case err != nil:
				for read != write {
					observer.Next(buffer[read])
					read = (read + 1) % n
				}
				observer.Error(err)
			case complete:
				for read != write {
					observer.Next(buffer[read])
					read = (read + 1) % n
				}
				observer.Complete()
			default:
				buffer[write] = next
				write = (write + 1) % n
				if write == read {
					read = (read + 1) % n
				}
			}
		}
		return filter
	}
	return factory
}

func (FiltersNamespace) IgnoreElements() FilterFactory {
	factory := func(observer Observer) Observer {
		filter := func(next interface{}, err error, complete bool) {
			if err != nil || complete {
				observer(next, err, complete)
			}
		}
		return filter
	}
	return factory
}

func (FiltersNamespace) One() FilterFactory {
	factory := func(observer Observer) Observer {
		count := 0
		var value interface{}
		filter := func(next interface{}, err error, complete bool) {
			if count > 1 {
				return
			}
			switch {
			case err != nil:
				observer.Error(err)
			case complete:
				if count == 2 {
					observer.Error(errors.New("expected one value"))
				} else {
					observer.Next(value)
					observer.Complete()
				}
			default:
				count++
				if count == 2 {
					observer.Error(errors.New("expected one value"))
					return
				}
				value = next
			}
		}
		return filter
	}
	return factory
}

type timedEntry struct {
	v interface{}
	t time.Time
}

func (FiltersNamespace) Replay(size int, duration time.Duration) FilterFactory {
	read := 0
	write := 0
	if size == 0 {
		size = MaxReplaySize
	}
	if duration == 0 {
		duration = time.Hour * 24 * 7 * 52
	}
	size++
	buffer := make([]timedEntry, size)
	factory := func(observer Observer) Observer {
		filter := func(next interface{}, err error, complete bool) {
			now := time.Now()
			switch {
			case err != nil:
				cursor := read
				for cursor != write {
					if buffer[cursor].t.After(now) {
						observer.Next(buffer[cursor].v)
					}
					cursor = (cursor + 1) % size
				}
				observer.Error(err)
			case complete:
				cursor := read
				for cursor != write {
					if buffer[cursor].t.After(now) {
						observer.Next(buffer[cursor].v)
					}
					cursor = (cursor + 1) % size
				}
				observer.Complete()
			default:
				buffer[write] = timedEntry{next, time.Now().Add(duration)}
				write = (write + 1) % size
				if write == read {
					if buffer[read].t.After(now) {
						observer.Next(buffer[read].v)
					}
					read = (read + 1) % size
				}
			}
		}
		return filter
	}
	return factory
}

func (FiltersNamespace) Sample(window time.Duration) FilterFactory {
	factory := func(observer Observer) Observer {
		mutex := &sync.Mutex{}
		cancel := make(chan bool, 1)
		var last interface{}
		haveNew := false
		go func() {
			for {
				select {
				case <-time.After(window):
					mutex.Lock()
					if haveNew {
						observer.Next(last)
						haveNew = false
					}
					mutex.Unlock()
				case <-cancel:
					return
				}
			}
		}()
		filter := func(next interface{}, err error, complete bool) {
			switch {
			case err != nil:
				cancel <- true
				observer.Error(err)
			case complete:
				cancel <- true
				observer.Complete()
			default:
				mutex.Lock()
				last = next
				haveNew = true
				mutex.Unlock()
			}
		}
		return filter
	}
	return factory
}

func (FiltersNamespace) Debounce(duration time.Duration) FilterFactory {
	factory := func(observer Observer) Observer {
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

type ObservableInt interface {
	Subscribe(IntObserver) Subscription
}

type IntObserver func(int, error, bool)

var zeroInt = *new(int)

func (f IntObserver) Next(next int) {
	f(next, nil, false)
}

func (f IntObserver) Error(err error) {
	f(zeroInt, err, false)
}

func (f IntObserver) Complete() {
	f(zeroInt, nil, true)
}

type IntSubscriber interface {
	Next(int)
	Error(error)
	Complete()
	Subscription
}

type IntCreateCallback func(IntSubscriber)

func (callback IntCreateCallback) Subscribe(observer IntObserver) Subscription {
	subscriber := &struct {
		IntObserver
		Subscription
	}{observer, new(Int32Subscription)}
	go callback(subscriber)
	return subscriber
}

type IntStream func(IntObserver) Subscription

func (s IntStream) Subscribe(observer IntObserver) Subscription {
	return s(observer)
}

func (s IntStream) SubscribeNext(f func(v int)) Subscription {
	return s.Subscribe(func(next int, err error, complete bool) {
		if err == nil && !complete {
			f(next)
		}
	})
}

// Wait for completion of the stream and return any error.
func (s IntStream) Wait() error {
	errch := make(chan error)
	s.Subscribe(func(next int, err error, complete bool) {
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

/////////////////////////////////////////////////////////////////////////////
// FROM
/////////////////////////////////////////////////////////////////////////////

// CreateInt calls f(subscriber) to produce values for a stream of ints.
func CreateInt(f IntCreateCallback) IntStream {
	return f.Subscribe
}

func Range(start, count int) IntStream {
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

func Interval(interval time.Duration) IntStream {
	return CreateInt(func(subscriber IntSubscriber) {
		i := 0
		for {
			time.Sleep(interval)
			if subscriber.Disposed() {
				return
			}
			subscriber.Next(i)
			i++
		}
	})
}

// Repeat value count times.
func RepeatInt(value, count int) IntStream {
	return CreateInt(func(subscriber IntSubscriber) {
		for i := 0; i < count; i++ {
			if subscriber.Disposed() {
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
// If the error is non-nil the returned IntStream will be that error,
// otherwise it will be a single-value stream of int.
func StartInt(f func() (int, error)) IntStream {
	return CreateInt(func(subscriber IntSubscriber) {
		if v, err := f(); err != nil {
			subscriber.Error(err)
		} else {
			subscriber.Next(v)
			subscriber.Complete()
		}
	})
}

func NeverInt() IntStream {
	return CreateInt(func(subscriber IntSubscriber) {})
}

func EmptyInt() IntStream {
	return CreateInt(func(subscriber IntSubscriber) {
		subscriber.Complete()
	})
}

func ThrowInt(err error) IntStream {
	return CreateInt(func(subscriber IntSubscriber) {
		subscriber.Error(err)
	})
}

func FromIntArray(array []int) IntStream {
	return CreateInt(func(subscriber IntSubscriber) {
		for _, v := range array {
			if subscriber.Disposed() {
				return
			}
			subscriber.Next(v)
		}
		subscriber.Complete()
		subscriber.Dispose()
	})
}

func FromInts(array ...int) IntStream {
	return FromIntArray(array)
}

func JustInt(element int) IntStream {
	return FromIntArray([]int{element})
}

func MergeInt(observables ...ObservableInt) IntStream {
	if len(observables) == 0 {
		return EmptyInt()
	}
	return (IntStream(observables[0].Subscribe)).Merge(observables[1:]...)
}

func MergeIntDelayError(observables ...ObservableInt) IntStream {
	if len(observables) == 0 {
		return EmptyInt()
	}
	return (IntStream(observables[0].Subscribe)).MergeDelayError(observables[1:]...)
}

func FromIntChannel(ch <-chan int) IntStream {
	return CreateInt(func(subscriber IntSubscriber) {
		for v := range ch {
			if subscriber.Disposed() {
				return
			}
			subscriber.Next(v)
		}
		subscriber.Complete()
	})
}

/////////////////////////////////////////////////////////////////////////////
// TO
/////////////////////////////////////////////////////////////////////////////

// ToOneWithError blocks until the stream emits exactly one value. Otherwise, it errors.
func (s IntStream) ToOneWithError() (int, error) {
	valuech := make(chan int, 1)
	errch := make(chan error, 1)
	s.One().Subscribe(func(next int, err error, complete bool) {
		if err != nil {
			errch <- err
		} else if !complete {
			valuech <- next
		}
	})
	select {
	case value := <-valuech:
		return value, nil
	case err := <-errch:
		return zeroInt, err
	}
}

// ToOne blocks and returns the only value emitted by the stream, or the zero
// value if an error occurs.
func (s IntStream) ToOne() int {
	value, _ := s.ToOneWithError()
	return value
}

// ToArrayWithError collects all values from the stream into an array,
// returning it and any error.
func (s IntStream) ToArrayWithError() ([]int, error) {
	array := []int{}
	completech := make(chan bool, 1)
	errch := make(chan error, 1)
	s.Subscribe(func(next int, err error, complete bool) {
		switch {
		case err != nil:
			errch <- err
		case complete:
			completech <- true
		default:
			array = append(array, next)
		}
	})
	select {
	case <-completech:
		return array, nil
	case err := <-errch:
		return array, err
	}
}

// ToArray blocks and returns the values from the stream in an array.
func (s IntStream) ToArray() []int {
	out, _ := s.ToArrayWithError()
	return out
}

// ToChannelWithError returns value and error channels corresponding to the stream elements and any error.
func (s IntStream) ToChannelWithError() (<-chan int, <-chan error) {
	ch := make(chan int, 1)
	errch := make(chan error, 1)
	s.Subscribe(func(next int, err error, complete bool) {
		switch {
		case err != nil:
			errch <- err
			close(errch)
			close(ch)
		case complete:
			close(ch)
		default:
			ch <- next
		}
	})
	return ch, errch
}

func (s IntStream) ToChannel() <-chan int {
	ch, _ := s.ToChannelWithError()
	return ch
}

/////////////////////////////////////////////////////////////////////////////
// FILTERS
/////////////////////////////////////////////////////////////////////////////

func (makeGenericFilter FilterFactory) FilterIntStream(source IntStream) IntStream {

	subscribe := func(sink IntObserver) Subscription {
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
		return source.Subscribe(int2Generic)
	}

	return subscribe
}

// Distinct removes duplicate elements in the stream.
func (s IntStream) Distinct() IntStream {
	return filters.Distinct().FilterIntStream(s)
}

// ElementAt yields the Nth element of the stream.
func (s IntStream) ElementAt(n int) IntStream {
	return filters.ElementAt(n).FilterIntStream(s)
}

// Filter elements in the stream on a function.
func (s IntStream) Filter(f func(int) bool) IntStream {
	filter := func(v interface{}) bool {
		return f(v.(int))
	}
	return filters.Filter(filter).FilterIntStream(s)
}

// Last returns just the first element of the stream.
func (s IntStream) First() IntStream {
	return filters.First().FilterIntStream(s)
}

// Last returns just the last element of the stream.
func (s IntStream) Last() IntStream {
	return filters.Last().FilterIntStream(s)
}

// SkipLast skips the first N elements of the stream.
func (s IntStream) Skip(n int) IntStream {
	return filters.Skip(n).FilterIntStream(s)
}

// SkipLast skips the last N elements of the stream.
func (s IntStream) SkipLast(n int) IntStream {
	return filters.SkipLast(n).FilterIntStream(s)
}

// Take returns just the first N elements of the stream.
func (s IntStream) Take(n int) IntStream {
	return filters.Take(n).FilterIntStream(s)
}

// TakeLast returns just the last N elements of the stream.
func (s IntStream) TakeLast(n int) IntStream {
	return filters.TakeLast(n).FilterIntStream(s)
}

// IgnoreElements ignores elements of the stream and emits only the completion events.
func (s IntStream) IgnoreElements() IntStream {
	return filters.IgnoreElements().FilterIntStream(s)
}

// IgnoreCompletion ignores the completion event of the stream and therefore returns a stream that never completes.
func (s IntStream) IgnoreCompletion() IntStream {
	subscribe := func(observer IntObserver) Subscription {
		filter := func(next int, err error, complete bool) {
			if !complete {
				observer(next, err, complete)
			}
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

func (s IntStream) One() IntStream {
	return filters.One().FilterIntStream(s)
}

func (s IntStream) Replay(size int, duration time.Duration) IntStream {
	return filters.Replay(size, duration).FilterIntStream(s)
}

func (s IntStream) Sample(duration time.Duration) IntStream {
	return filters.Sample(duration).FilterIntStream(s)
}

// Debounce reduces subsequent duplicates to single items during a certain duration
func (s IntStream) Debounce(duration time.Duration) IntStream {
	return filters.Debounce(duration).FilterIntStream(s)
}

/////////////////////////////////////////////////////////////////////////////
// COUNT
/////////////////////////////////////////////////////////////////////////////

func (s IntStream) Count() IntStream {
	subscribe := func(observer IntObserver) Subscription {
		count := 0
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
		return s.Subscribe(filter)
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// CONCAT
/////////////////////////////////////////////////////////////////////////////

type concatInt struct {
	observables []ObservableInt
}

func (m *concatInt) Subscribe(observer IntObserver) Subscription {
	if len(m.observables) == 0 {
		observer.Complete()
		return ClosedSubscription{}
	}
	type switchingIntObserver struct {
		switchFunc   IntObserver
		count        int
		subscription Subscription
	}

	var obs *switchingIntObserver
	switchFunc := func(next int, err error, complete bool) {
		switch {
		case err != nil:
			observer.Error(err)
			obs.count = len(m.observables)
			obs.subscription.Dispose()
		case complete:
			obs.count++
			if obs.count >= len(m.observables) {
				observer.Complete()
				obs.subscription.Dispose()
			} else {
				obs.subscription = m.observables[obs.count].Subscribe(obs.switchFunc)
			}
		default:
			observer.Next(next)
		}
	}
	obs = &switchingIntObserver{switchFunc, 0, nil}
	obs.subscription = m.observables[0].Subscribe(obs.switchFunc)
	return new(Int32Subscription)
}

func (s IntStream) Concat(observables ...ObservableInt) IntStream {
	return (&concatInt{append([]ObservableInt{s}, observables...)}).Subscribe
}

/////////////////////////////////////////////////////////////////////////////
// MERGE
/////////////////////////////////////////////////////////////////////////////

type mergeInt struct {
	delayError  bool
	observables []ObservableInt
}

func (m *mergeInt) Subscribe(observer IntObserver) Subscription {
	lock := sync.Mutex{}
	completed := 0
	var firstError error
	relay := func(next int, err error, complete bool) {
		lock.Lock()
		defer lock.Unlock()
		if completed >= len(m.observables) {
			return
		}

		switch {
		case err != nil:
			if m.delayError {
				firstError = err
				completed++
			} else {
				observer.Error(err)
				completed = len(m.observables)
			}

		case complete:
			completed++
			if completed == len(m.observables) {
				if firstError != nil {
					observer.Error(firstError)
				} else {
					observer.Complete()
				}
			}
		default:
			observer.Next(next)
		}
	}
	for _, observable := range m.observables {
		observable.Subscribe(relay)
	}
	return new(Int32Subscription)
}

// Merge an arbitrary number of observables with this one.
// An error from any of the observables will terminate the merged stream.
func (s IntStream) Merge(other ...ObservableInt) IntStream {
	if len(other) == 0 {
		return s
	}
	return (&mergeInt{false, append(other, s)}).Subscribe
}

// Merge an arbitrary number of observables with this one.
// Any error will be deferred until all observables terminate.
func (s IntStream) MergeDelayError(other ...ObservableInt) IntStream {
	if len(other) == 0 {
		return s
	}
	return (&mergeInt{true, append(other, s)}).Subscribe
}

/////////////////////////////////////////////////////////////////////////////
// CATCH
/////////////////////////////////////////////////////////////////////////////

type catchInt struct {
	parent ObservableInt
	catch  ObservableInt
}

func (r *catchInt) Subscribe(observer IntObserver) Subscription {
	subscription := new(Int32Subscription)
	run := func(next int, err error, complete bool) {
		switch {
		case err != nil:
			r.catch.Subscribe(observer)
		case complete:
			observer.Complete()
		default:
			observer.Next(next)
		}
	}
	r.parent.Subscribe(run)
	return subscription
}

func (s IntStream) Catch(catch ObservableInt) IntStream {
	return (&catchInt{s, catch}).Subscribe
}

/////////////////////////////////////////////////////////////////////////////
// RETRY
/////////////////////////////////////////////////////////////////////////////

type retryIntObserver struct {
	observable ObservableInt
	observer   IntObserver
}

func (r *retryIntObserver) retry(next int, err error, complete bool) {
	switch {
	case err != nil:
		r.observable.Subscribe(IntObserver(r.retry))
	case complete:
		r.observer.Complete()
	default:
		r.observer.Next(next)
	}
}

type retryInt struct {
	observable ObservableInt
}

func (r *retryInt) Subscribe(observer IntObserver) Subscription {
	ro := &retryIntObserver{r.observable, observer}
	r.observable.Subscribe(IntObserver(ro.retry))
	return new(Int32Subscription)
}

func (s IntStream) Retry() IntStream {
	return (&retryInt{s}).Subscribe
}

/////////////////////////////////////////////////////////////////////////////
// DO
/////////////////////////////////////////////////////////////////////////////

// Do applies a function for each value passing through the stream.
func (s IntStream) Do(f func(next int)) IntStream {
	subscribe := func(observer IntObserver) Subscription {
		filter := func(next int, err error, complete bool) {
			if err == nil && !complete {
				f(next)
			}
			observer(next, err, complete)
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

// DoOnError applies a function for any error on the stream.
func (s IntStream) DoOnError(f func(err error)) IntStream {
	subscribe := func(observer IntObserver) Subscription {
		filter := func(next int, err error, complete bool) {
			if err != nil {
				f(err)
			}
			observer(next, err, complete)
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

// DoOnComplete applies a function when the stream completes.
func (s IntStream) DoOnComplete(f func()) IntStream {
	subscribe := func(observer IntObserver) Subscription {
		filter := func(next int, err error, complete bool) {
			if complete {
				f()
			}
			observer(next, err, complete)
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

// DoOnFinally applies a function for any error or completion on the stream, using err == nil to indicate completion.
func (s IntStream) DoOnFinally(f func(err error)) IntStream {
	subscribe := func(observer IntObserver) Subscription {
		filter := func(next int, err error, complete bool) {
			if err != nil || complete {
				f(err)
			}
			observer(next, err, complete)
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// REDUCE
/////////////////////////////////////////////////////////////////////////////

func (s IntStream) Reduce(initial int, reducer func(int, int) int) IntStream {
	value := initial
	subscribe := func(observer IntObserver) Subscription {
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
		return s.Subscribe(filter)
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// SCAN
/////////////////////////////////////////////////////////////////////////////

func (s IntStream) Scan(initial int, f func(int, int) int) IntStream {
	value := initial
	subscribe := func(observer IntObserver) Subscription {
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
		return s.Subscribe(filter)
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// TIMEOUT
/////////////////////////////////////////////////////////////////////////////

type timeoutInt struct {
	parent  ObservableInt
	timeout time.Duration
}

func (t *timeoutInt) Subscribe(observer IntObserver) Subscription {
	subscription := NewChannelSubscription()
	cancel := t.parent.Subscribe(observer)
	go func() {
		select {
		case <-time.After(t.timeout):
			observer.Error(ErrTimeout)
			cancel.Dispose()
			subscription.Dispose()
		case <-subscription:
			cancel.Dispose()
		}
	}()
	return subscription
}

func (s IntStream) Timeout(timeout time.Duration) IntStream {
	return (&timeoutInt{s, timeout}).Subscribe
}

/////////////////////////////////////////////////////////////////////////////
// FORK
/////////////////////////////////////////////////////////////////////////////

type forkedIntStream struct {
	lock      sync.Mutex
	parent    ObservableInt
	observers []IntObserver
}

func (f *forkedIntStream) Subscribe(observer IntObserver) Subscription {
	f.lock.Lock()
	defer f.lock.Unlock()
	i := len(f.observers)
	f.observers = append(f.observers, observer)
	sub := new(CallbackSubscription)
	*sub = CallbackSubscription(func() {
		f.lock.Lock()
		defer f.lock.Unlock()
		f.observers[i] = nil
	})
	return sub
}

// Fork replicates each event from the parent to every subscriber of the fork.
func (s IntStream) Fork() IntStream {
	f := &forkedIntStream{parent: s}
	s.Subscribe(IntObserver(func(n int, err error, complete bool) {
		f.lock.Lock()
		defer f.lock.Unlock()
		for _, o := range f.observers {
			if o == nil {
				continue
			}
			switch {
			case err != nil:
				o.Error(err)
			case complete:
				o.Complete()
			default:
				o.Next(n)
			}
		}
	}))
	return f.Subscribe
}

/////////////////////////////////////////////////////////////////////////////
// MATHEMATICAL
/////////////////////////////////////////////////////////////////////////////

func (s IntStream) Average() IntStream {
	var sum int
	var count int
	subscribe := func(observer IntObserver) Subscription {
		filter := func(next int, err error, complete bool) {
			switch {
			case err != nil:
				observer.Next(sum / count)
				observer.Error(err)
			case complete:
				observer.Next(sum / count)
				observer.Complete()
			default:
				sum += next
				count++
			}
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

func (s IntStream) Sum() IntStream {
	var sum int
	subscribe := func(observer IntObserver) Subscription {
		filter := func(next int, err error, complete bool) {
			switch {
			case err != nil:
				observer.Next(sum)
				observer.Error(err)
			case complete:
				observer.Next(sum)
				observer.Complete()
			default:
				sum += next
			}
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

func (s IntStream) Min() IntStream {
	started := false
	var min int
	subscribe := func(observer IntObserver) Subscription {
		filter := func(next int, err error, complete bool) {
			switch {
			case err != nil:
				if started {
					observer.Next(min)
				}
				observer.Error(err)
			case complete:
				if started {
					observer.Next(min)
				}
				observer.Complete()
			default:
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
		return s.Subscribe(filter)
	}
	return subscribe
}

func (s IntStream) Max() IntStream {
	started := false
	var max int
	subscribe := func(observer IntObserver) Subscription {
		filter := func(next int, err error, complete bool) {
			switch {
			case err != nil:
				if started {
					observer.Next(max)
				}
				observer.Error(err)
			case complete:
				if started {
					observer.Next(max)
				}
				observer.Complete()
			default:
				if started {
					if max <= next {
						max = next
					}
				} else {
					max = next
					started = true
				}
			}
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// MAP and FLATMAP (int -> int)
/////////////////////////////////////////////////////////////////////////////

func (s IntStream) Map(f func(int) int) IntStream {
	subscribe := func(observer IntObserver) Subscription {
		filter := func(next int, err error, complete bool) {
			var mapped int
			if err == nil && !complete {
				mapped = f(next)
			}
			observer(mapped, err, complete)
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

func (s IntStream) FlatMap(f func(int) IntStream) IntStream {

	subscribe := func(observer IntObserver) Subscription {
		subscription := new(Int32Subscription)
		var wg sync.WaitGroup
		var mx sync.Mutex
		gatedObserver := func(next int, err error, complete bool) {
			mx.Lock()
			defer mx.Unlock()
			observer(next, err, complete)
		}
		filter := func(next int, err error, complete bool) {
			switch {
			case err != nil:
				wg.Wait()
				observer.Error(err)
			case complete:
				wg.Wait()
				observer.Complete()
			default:
				wg.Add(1)
				nextStream := f(next)
				nextStream.DoOnFinally(func(error) { wg.Done() }).IgnoreCompletion().Subscribe(gatedObserver)
			}
		}
		s.Subscribe(filter)
		return subscription
	}

	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// MAP and FLATMAP (int -> float64)
/////////////////////////////////////////////////////////////////////////////

func (s IntStream) MapFloat64(f func(int) float64) Float64Stream {
	subscribe := func(observer Float64Observer) Subscription {
		filter := func(next int, err error, complete bool) {
			var mapped float64
			if err == nil && !complete {
				mapped = f(next)
			}
			observer(mapped, err, complete)
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

func (s IntStream) FlatMapFloat64(f func(int) Float64Stream) Float64Stream {

	subscribe := func(observer Float64Observer) Subscription {
		var wg sync.WaitGroup
		var mx sync.Mutex
		gatedObserver := func(next float64, err error, complete bool) {
			mx.Lock()
			defer mx.Unlock()
			observer(next, err, complete)
		}
		filter := func(next int, err error, complete bool) {
			switch {
			case err != nil:
				wg.Wait()
				observer.Error(err)
			case complete:
				wg.Wait()
				observer.Complete()
			default:
				wg.Add(1)
				nextStream := f(next)
				nextStream.DoOnFinally(func(error) { wg.Done() }).IgnoreCompletion().Subscribe(gatedObserver)
			}
		}
		s.Subscribe(filter)
		return new(Int32Subscription)
	}

	return subscribe
}

/////////////////////////////////////////////////////////////////////////////
// MAP and FLATMAP (int -> string)
/////////////////////////////////////////////////////////////////////////////

func (s IntStream) MapString(f func(int) string) StringStream {
	subscribe := func(observer StringObserver) Subscription {
		filter := func(next int, err error, complete bool) {
			var mapped string
			if err == nil && !complete {
				mapped = f(next)
			}
			observer(mapped, err, complete)
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

func (s IntStream) FlatMapString(f func(int) StringStream) StringStream {

	subscribe := func(observer StringObserver) Subscription {
		var wg sync.WaitGroup
		var mx sync.Mutex
		gatedObserver := func(next string, err error, complete bool) {
			mx.Lock()
			defer mx.Unlock()
			observer(next, err, complete)
		}
		filter := func(next int, err error, complete bool) {
			switch {
			case err != nil:
				wg.Wait()
				observer.Error(err)
			case complete:
				wg.Wait()
				observer.Complete()
			default:
				wg.Add(1)
				nextStream := f(next)
				nextStream.DoOnFinally(func(error) { wg.Done() }).IgnoreCompletion().Subscribe(gatedObserver)
			}
		}
		s.Subscribe(filter)
		return new(Int32Subscription)
	}

	return subscribe
}

////////////////////////////////////////////////////////
// Float64
////////////////////////////////////////////////////////

type ObservableFloat64 interface {
	Subscribe(Float64Observer) Subscription
}

type Float64Observer func(float64, error, bool)

var zeroFloat64 = *new(float64)

func (f Float64Observer) Next(next float64) {
	f(next, nil, false)
}

func (f Float64Observer) Error(err error) {
	f(zeroFloat64, err, false)
}

func (f Float64Observer) Complete() {
	f(zeroFloat64, nil, true)
}

type Float64Stream func(Float64Observer) Subscription

func (s Float64Stream) Subscribe(observer Float64Observer) Subscription {
	return s(observer)
}

// IgnoreCompletion ignores the completion event of the stream and therefore returns a stream that never completes.
func (s Float64Stream) IgnoreCompletion() Float64Stream {
	subscribe := func(observer Float64Observer) Subscription {
		filter := func(next float64, err error, complete bool) {
			if !complete {
				observer(next, err, complete)
			}
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

// DoOnFinally applies a function for any error or completion on the stream, using err == nil to indicate completion.
func (s Float64Stream) DoOnFinally(f func(err error)) Float64Stream {
	subscribe := func(observer Float64Observer) Subscription {
		filter := func(next float64, err error, complete bool) {
			if err != nil || complete {
				f(err)
			}
			observer(next, err, complete)
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

////////////////////////////////////////////////////////
// String
////////////////////////////////////////////////////////

type ObservableString interface {
	Subscribe(StringObserver) Subscription
}

type StringObserver func(string, error, bool)

var zeroString = *new(string)

func (f StringObserver) Next(next string) {
	f(next, nil, false)
}

func (f StringObserver) Error(err error) {
	f(zeroString, err, false)
}

func (f StringObserver) Complete() {
	f(zeroString, nil, true)
}

type StringStream func(StringObserver) Subscription

func (s StringStream) Subscribe(observer StringObserver) Subscription {
	return s(observer)
}

// ToArrayWithError collects all values from the stream stringo an array,
// returning it and any error.
func (s StringStream) ToArrayWithError() ([]string, error) {
	array := []string{}
	completech := make(chan bool, 1)
	errch := make(chan error, 1)
	s.Subscribe(func(next string, err error, complete bool) {
		switch {
		case err != nil:
			errch <- err
		case complete:
			completech <- true
		default:
			array = append(array, next)
		}
	})
	select {
	case <-completech:
		return array, nil
	case err := <-errch:
		return array, err
	}
}

// ToArray blocks and returns the values from the stream in an array.
func (s StringStream) ToArray() []string {
	out, _ := s.ToArrayWithError()
	return out
}

// IgnoreCompletion ignores the completion event of the stream and therefore returns a stream that never completes.
func (s StringStream) IgnoreCompletion() StringStream {
	subscribe := func(observer StringObserver) Subscription {
		filter := func(next string, err error, complete bool) {
			if !complete {
				observer(next, err, complete)
			}
		}
		return s.Subscribe(filter)
	}
	return subscribe
}

// DoOnFinally applies a function for any error or completion on the stream, using err == nil to indicate completion.
func (s StringStream) DoOnFinally(f func(err error)) StringStream {
	subscribe := func(observer StringObserver) Subscription {
		filter := func(next string, err error, complete bool) {
			if err != nil || complete {
				f(err)
			}
			observer(next, err, complete)
		}
		return s.Subscribe(filter)
	}
	return subscribe
}
