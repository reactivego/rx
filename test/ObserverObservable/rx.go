// Code generated by jig; DO NOT EDIT.

//go:generate jig

package ObserverObservable

import (
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/reactivego/scheduler"
	"github.com/reactivego/rx/subscriber"
)

//jig:name Observer

// Observer is a function that gets called whenever the Observable has
// something to report. The next argument is the item value that is only
// valid when the done argument is false. When done is true and the err
// argument is not nil, then the Observable has terminated with an error.
// When done is true and the err argument is nil, then the Observable has
// completed normally.
type Observer func(next interface{}, err error, done bool)

//jig:name Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler = scheduler.Scheduler

//jig:name Subscriber

// Subscriber is an interface that can be passed in when subscribing to an
// Observable. It allows a set of observable subscriptions to be canceled
// from a single subscriber at the root of the subscription tree.
type Subscriber = subscriber.Subscriber

//jig:name Observable

// Observable is a function taking an Observer, Scheduler and Subscriber.
// Calling it will subscribe the Observer to events from the Observable.
type Observable func(Observer, Scheduler, Subscriber)

//jig:name MakeObserverObservable

const OutOfSubscriptions = RxError("out of subscriptions")

// MakeObserverObservable turns an observer into a multicasting and buffering
// observable. Both the observer and the obeservable are returned. These are
// then used as the core of any Subject implementation. The Observer side is
// used to pass items into the buffering multicaster. This then multicasts the
// items to every Observer that subscribes to the returned Observable.
//
//	age     age below which items are kept to replay to a new subscriber.
//	length  length of the item buffer, number of items kept to replay to a new subscriber.
//	[cap]   Capacity of the item buffer, number of items that can be observed before blocking.
//	[scap]  Capacity of the subscription list, max number of simultaneous subscribers.
func MakeObserverObservable(age time.Duration, length int, capacity ...int) (Observer, Observable) {
	const (
		ms	= time.Millisecond
		us	= time.Microsecond
	)

	type subscription struct {
		cursor		uint64
		state		uint64		// active, canceled, closed
		activated	time.Time	// track activity to deterime backoff
		subscribeOn	Scheduler
	}

	// cursor
	const (
		maxuint64 uint64 = math.MaxUint64	// park unused cursor at maxuint64
	)

	// state
	const (
		active	uint64	= iota
		canceled
		closing
		closed
	)

	type subscriptions struct {
		sync.Mutex
		*sync.Cond
		entries	[]subscription
		access	uint32	// unlocked, locked
	}

	// access
	const (
		unlocked	uint32	= iota
		locked
	)

	type item struct {
		Value	interface{}
		At	time.Time
	}

	type buffer struct {
		age	time.Duration
		keep	uint64
		mod	uint64
		size	uint64

		items	[]item
		begin	uint64
		end	uint64
		commit	uint64
		state	uint64	// active, closed

		subscriptions	subscriptions

		err	error
	}

	make := func(age time.Duration, length int, capacity ...int) *buffer {
		if length < 0 {
			length = 0
		}
		keep := uint64(length)

		cap, scap := length, 32
		switch {
		case len(capacity) >= 2:
			cap, scap = capacity[0], capacity[1]
		case len(capacity) == 1:
			cap = capacity[0]
		}
		if cap < length {
			cap = length
		}
		if cap == 0 {
			cap = 1
		}
		size := uint64(1) << uint(math.Ceil(math.Log2(float64(cap))))

		if scap < 1 {
			scap = 1
		}
		buf := &buffer{
			age:	age,
			keep:	keep,
			mod:	size - 1,
			size:	size,

			items:	make([]item, size),
			end:	size,
			subscriptions: subscriptions{
				entries: make([]subscription, 0, scap),
			},
		}
		buf.subscriptions.Cond = sync.NewCond(&buf.subscriptions.Mutex)
		return buf
	}
	buf := make(age, length, capacity...)

	accessSubscriptions := func(access func([]subscription)) bool {
		gosched := false
		for !atomic.CompareAndSwapUint32(&buf.subscriptions.access, unlocked, locked) {
			runtime.Gosched()
			gosched = true
		}
		access(buf.subscriptions.entries)
		atomic.StoreUint32(&buf.subscriptions.access, unlocked)
		return gosched
	}

	send := func(value interface{}) {
		for buf.commit == buf.end {
			full := false
			subscribeOn := Scheduler(nil)
			gosched := accessSubscriptions(func(subscriptions []subscription) {
				slowest := maxuint64
				for i := range subscriptions {
					current := atomic.LoadUint64(&subscriptions[i].cursor)
					if current < slowest {
						slowest = current
						subscribeOn = subscriptions[i].subscribeOn
					}
				}
				end := atomic.LoadUint64(&buf.end)
				if atomic.LoadUint64(&buf.begin) < slowest && slowest <= end {
					if slowest+buf.keep > end {
						slowest = end - buf.keep + 1
					}
					atomic.StoreUint64(&buf.begin, slowest)
					atomic.StoreUint64(&buf.end, slowest+buf.size)
				} else {
					if slowest == maxuint64 {
						atomic.AddUint64(&buf.begin, 1)
						atomic.AddUint64(&buf.end, 1)
					} else {
						full = true
					}
				}
			})
			if full {
				if !gosched {
					if subscribeOn != nil {
						subscribeOn.Gosched()
					} else {
						runtime.Gosched()
					}
				}
				if atomic.LoadUint64(&buf.state) != active {
					return
				}
			}
		}
		buf.items[buf.commit&buf.mod] = item{Value: value, At: time.Now()}
		atomic.AddUint64(&buf.commit, 1)
		buf.subscriptions.Broadcast()
	}

	close := func(err error) {
		if atomic.CompareAndSwapUint64(&buf.state, active, closing) {
			buf.err = err
			if atomic.CompareAndSwapUint64(&buf.state, closing, closed) {
				accessSubscriptions(func(subscriptions []subscription) {
					for i := range subscriptions {
						atomic.CompareAndSwapUint64(&subscriptions[i].state, active, closed)
					}
				})
			}
		}
		buf.subscriptions.Broadcast()
	}

	observer := func(next interface{}, err error, done bool) {
		if atomic.LoadUint64(&buf.state) == active {
			if !done {
				send(next)
			} else {
				close(err)
			}
		}
	}

	appendSubscription := func(subscribeOn Scheduler) (sub *subscription, err error) {
		accessSubscriptions(func([]subscription) {
			cursor := atomic.LoadUint64(&buf.begin)
			s := &buf.subscriptions
			if len(s.entries) < cap(s.entries) {
				s.entries = append(s.entries, subscription{cursor: cursor, subscribeOn: subscribeOn})
				sub = &s.entries[len(s.entries)-1]
				return
			}
			for i := range s.entries {
				sub = &s.entries[i]
				if atomic.CompareAndSwapUint64(&sub.cursor, maxuint64, cursor) {
					sub.subscribeOn = subscribeOn
					return
				}
			}
			sub = nil
			err = OutOfSubscriptions
			return
		})
		return
	}

	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		sub, err := appendSubscription(subscribeOn)
		if err != nil {
			runner := subscribeOn.Schedule(func() {
				if subscriber.Subscribed() {
					observe(nil, err, true)
				}
			})
			subscriber.OnUnsubscribe(runner.Cancel)
			return
		}
		commit := atomic.LoadUint64(&buf.commit)
		if atomic.LoadUint64(&buf.begin)+buf.keep < commit {
			atomic.StoreUint64(&sub.cursor, commit-buf.keep)
		}
		atomic.StoreUint64(&sub.state, atomic.LoadUint64(&buf.state))
		sub.activated = time.Now()

		receiver := subscribeOn.ScheduleFutureRecursive(0, func(self func(time.Duration)) {
			commit := atomic.LoadUint64(&buf.commit)

			if sub.cursor == commit {
				if atomic.CompareAndSwapUint64(&sub.state, canceled, canceled) {

					atomic.StoreUint64(&sub.cursor, maxuint64)
					return
				} else {

					now := time.Now()
					if now.Before(sub.activated.Add(1 * ms)) {

						self(50 * us)
						return
					} else if now.Before(sub.activated.Add(250 * ms)) {
						if atomic.CompareAndSwapUint64(&sub.state, closed, closed) {

							observe(nil, buf.err, true)
							atomic.StoreUint64(&sub.cursor, maxuint64)
							return
						}

						self(500 * us)
						return
					} else {
						if subscribeOn.IsConcurrent() {

							buf.subscriptions.Lock()
							buf.subscriptions.Wait()
							buf.subscriptions.Unlock()
							sub.activated = time.Now()
							self(0)
							return
						} else {

							self(5 * ms)
							return
						}
					}
				}
			}

			if atomic.LoadUint64(&sub.state) == canceled {
				atomic.StoreUint64(&sub.cursor, maxuint64)
				return
			}
			for ; sub.cursor != commit; atomic.AddUint64(&sub.cursor, 1) {
				item := &buf.items[sub.cursor&buf.mod]
				if buf.age == 0 || item.At.IsZero() || time.Since(item.At) < buf.age {
					observe(item.Value, nil, false)
				}
				if atomic.LoadUint64(&sub.state) == canceled {
					atomic.StoreUint64(&sub.cursor, maxuint64)
					return
				}
			}

			sub.activated = time.Now()
			self(0)
		})
		subscriber.OnUnsubscribe(receiver.Cancel)

		subscriber.OnUnsubscribe(func() {
			atomic.CompareAndSwapUint64(&sub.state, active, canceled)
			buf.subscriptions.Broadcast()
		})
	}
	return observer, observable
}

//jig:name RxError

type RxError string

func (e RxError) Error() string	{ return string(e) }

//jig:name Observable_SubscribeOn

// SubscribeOn specifies the scheduler an Observable should use when it is
// subscribed to.
func (o Observable) SubscribeOn(scheduler Scheduler) Observable {
	observable := func(observe Observer, _ Scheduler, subscriber Subscriber) {
		if scheduler.IsConcurrent() {
			subscriber.OnWait(nil)
		} else {
			subscriber.OnWait(scheduler.Wait)
		}
		o(observe, scheduler, subscriber)
	}
	return observable
}

//jig:name Subscription

// Subscription is an alias for the subscriber.Subscription interface type.
type Subscription = subscriber.Subscription

//jig:name Observable_Subscribe

// Subscribe operates upon the emissions and notifications from an Observable.
// This method returns a Subscription.
// Subscribe uses a trampoline scheduler created with scheduler.MakeTrampoline().
func (o Observable) Subscribe(observe Observer, schedulers ...Scheduler) Subscription {
	subscriber := subscriber.New()
	schedulers = append(schedulers, scheduler.MakeTrampoline())
	observer := func(next interface{}, err error, done bool) {
		if !done {
			observe(next, err, done)
		} else {
			var zero interface{}
			observe(zero, err, true)
			subscriber.Done(err)
		}
	}
	if !schedulers[0].IsConcurrent() {
		subscriber.OnWait(schedulers[0].Wait)
	}
	o(observer, schedulers[0], subscriber)
	return subscriber
}
