package rx

import (
	"sync"
	"time"

	"github.com/reactivego/scheduler"
)

//jig:template Observable Audit

// Audit waits until the source emits and then starts a timer. When the timer
// expires, Audit will emit the last value received from the source during the
// time period when the timer was active.
func (o Observable) Audit(duration time.Duration) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var audit struct {
			sync.Mutex
			runner scheduler.Runner
			next   interface{}
			done   bool
		}
		auditer := func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				audit.Lock()
				audit.runner = nil
				next := audit.next
				done := audit.done
				audit.Unlock()
				if !done {
					observe(next, nil, false)
				}
			}
		}
		observer := func(next interface{}, err error, done bool) {
			if subscriber.Subscribed() {
				if !done {
					audit.Lock()
					audit.next = next
					if audit.runner == nil {
						audit.runner = subscribeOn.ScheduleFutureRecursive(duration, auditer)
					}
					audit.Unlock()
				} else {
					audit.Lock()
					audit.done = true
					audit.Unlock()
					observe(nil, err, true)
				}
			}
		}
		subscriber.OnUnsubscribe(func() {
			audit.Lock()
			if audit.runner != nil {
				audit.runner.Cancel()
			}
			audit.Unlock()
		})
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Audit
//jig:needs Observable Audit

// Audit waits until the source emits and then starts a timer. When the timer
// expires, Audit will emit the last value received from the source during the
// time period when the timer was active.
func (o ObservableFoo) Audit(duration time.Duration) ObservableFoo {
	return o.AsObservable().Audit(duration).AsObservableFoo()
}

//jig:template Observable Debounce

// Debounce only emits the last item of a burst from an Observable if a
// particular timespan has passed without it emitting another item.
func (o Observable) Debounce(duration time.Duration) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var debounce struct {
			sync.Mutex
			runner scheduler.Runner
			next   interface{}
			done   bool
		}
		debouncer := func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				debounce.Lock()
				debounce.runner = nil
				next := debounce.next
				done := debounce.done
				debounce.Unlock()
				if !done {
					observe(next, nil, false)
				}
			}
		}
		observer := func(next interface{}, err error, done bool) {
			if subscriber.Subscribed() {
				if !done {
					debounce.Lock()
					debounce.next = next
					if debounce.runner != nil {
						debounce.runner.Cancel()
					}
					debounce.runner = subscribeOn.ScheduleFutureRecursive(duration, debouncer)
					debounce.Unlock()
				} else {
					debounce.Lock()
					debounce.done = true
					debounce.Unlock()
					observe(nil, err, true)
				}
			}
		}
		subscriber.OnUnsubscribe(func() {
			debounce.Lock()
			if debounce.runner != nil {
				debounce.runner.Cancel()
			}
			debounce.Unlock()
		})
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Debounce
//jig:needs Observable Debounce

// Debounce only emits the last item of a burst from an ObservableFoo if a
// particular timespan has passed without it emitting another item.
func (o ObservableFoo) Debounce(duration time.Duration) ObservableFoo {
	return o.AsObservable().Debounce(duration).AsObservableFoo()
}

//jig:template Observable Delay

// Delay shifts an emission from an Observable forward in time by a particular
// amount of time. The relative time intervals between emissions are preserved.
func (o Observable) Delay(duration time.Duration) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		type emission struct {
			at   time.Time
			next interface{}
			err  error
			done bool
		}
		var delay struct {
			sync.Mutex
			emissions []emission
		}
		delayer := subscribeOn.ScheduleFutureRecursive(duration, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				delay.Lock()
				for _, entry := range delay.emissions {
					delay.Unlock()
					due := entry.at.Sub(subscribeOn.Now())
					if due > 0 {
						self(due)
						return
					}
					observe(entry.next, entry.err, entry.done)
					if entry.done || subscriber.Canceled() {
						return
					}
					delay.Lock()
					delay.emissions = delay.emissions[1:]
				}
				delay.Unlock()
				self(duration) // keep on rescheduling the emitter
			}
		})
		subscriber.OnUnsubscribe(delayer.Cancel)
		observer := func(next interface{}, err error, done bool) {
			delay.Lock()
			delay.emissions = append(delay.emissions, emission{subscribeOn.Now().Add(duration), next, err, done})
			delay.Unlock()
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Delay
//jig:needs Observable Delay

// Delay shifts an emission from an Observable forward in time by a particular
// amount of time. The relative time intervals between emissions are preserved.
func (o ObservableFoo) Delay(duration time.Duration) ObservableFoo {
	return o.AsObservable().Delay(duration).AsObservableFoo()
}

//jig:template Interval
//jig:needs ObservableInt

// Interval creates an ObservableInt that emits a sequence of integers spaced
// by a particular time interval. First integer is not emitted immediately, but
// only after the first time interval has passed.
func Interval(interval time.Duration) ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		runner := subscribeOn.ScheduleFutureRecursive(interval, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				observe(i, nil, false)
				i++
				if subscriber.Subscribed() {
					self(interval)
				}
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template Observable Sample

// Sample emits the most recent item emitted by an Observable within periodic time intervals.
func (o Observable) Sample(window time.Duration) Observable {
	observable := Observable(func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var sample struct {
			sync.Mutex
			at   time.Time
			next interface{}
			done bool
		}
		sampler := subscribeOn.ScheduleFutureRecursive(window, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				sample.Lock()
				if !sample.done {
					begin := subscribeOn.Now().Add(-window)
					if !sample.at.Before(begin) {
						observe(sample.next, nil, false)
					}
					if subscriber.Subscribed() {
						self(window)
					}
				}
				sample.Unlock()
			}
		})
		subscriber.OnUnsubscribe(sampler.Cancel)
		observer := func(next interface{}, err error, done bool) {
			if subscriber.Subscribed() {
				sample.Lock()
				sample.at = subscribeOn.Now()
				sample.next = next
				sample.done = done
				sample.Unlock()
				if done {
					observe(nil, err, true)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	})
	return observable
}

//jig:template Observable<Foo> Sample
//jig:needs Observable Sample

// Sample emits the most recent item emitted by an ObservableFoo within periodic
// time intervals.
func (o ObservableFoo) Sample(window time.Duration) ObservableFoo {
	return o.AsObservable().Sample(window).AsObservableFoo()
}

//jig:template Observable Throttle

// Throttle emits when the source emits and then starts a timer during which
// all emissions from the source are ignored. After the timer expires, Throttle
// will again emit the next item the source emits, and so on.
func (o Observable) Throttle(duration time.Duration) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var deadline time.Time
		observer := func(next interface{}, err error, done bool) {
			if !done {
				if subscribeOn.Now().After(deadline) {
					observe(next, nil, false)
					deadline = subscribeOn.Now().Add(duration)
				}
			} else {
				observe(nil, err, true)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Throttle
//jig:needs Observable Throttle

// Throttle
func (o ObservableFoo) Throttle(duration time.Duration) ObservableFoo {
	return o.AsObservable().Throttle(duration).AsObservableFoo()
}

//jig:type Time time.Time

//jig:template Time

type Time = time.Time

//jig:template Ticker
//jig:needs Time, ObservableTime

// Ticker creates an ObservableTime that emits a sequence of timestamps after
//an initialDelay has passed. Subsequent timestamps are emitted using a
// schedule of intervals passed in. If only the initialDelay is given, Tichker
// will emit only once.
func Ticker(initialDelay time.Duration, intervals ...time.Duration) ObservableTime {
	observable := func(observe TimeObserver, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		runner := subscribeOn.ScheduleFutureRecursive(initialDelay, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				observe(subscribeOn.Now(), nil, false)
				if subscriber.Subscribed() {
					if len(intervals) > 0 {
						self(intervals[i%len(intervals)])
					}
				}
				i++
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template Observable Timeout
//jig:needs RxError, Observable Serialize, nil

// ErrTimeout is delivered to an observer if the stream times out.
const ErrTimeout = RxError("timeout")

// Timeout mirrors the source Observable, but issues an error notification if a
// particular period of time elapses without any emitted items.
// Timeout schedules a task on the scheduler passed to it during subscription.
func (o Observable) Timeout(due time.Duration) Observable {
	observable := Observable(func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var timeout struct {
			sync.Mutex
			at       time.Time
			occurred bool
		}
		timeout.at = subscribeOn.Now().Add(due)
		timer := subscribeOn.ScheduleFutureRecursive(due, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				timeout.Lock()
				if !timeout.occurred {
					due := timeout.at.Sub(subscribeOn.Now())
					if due > 0 {
						self(due)
					} else {
						timeout.occurred = true
						timeout.Unlock()
						observe(nil, ErrTimeout, true)
						timeout.Lock()
					}
				}
				timeout.Unlock()
			}
		})
		subscriber.OnUnsubscribe(timer.Cancel)
		observer := func(next interface{}, err error, done bool) {
			if subscriber.Subscribed() {
				timeout.Lock()
				if !timeout.occurred {
					now := subscribeOn.Now()
					if now.Before(timeout.at) {
						timeout.at = now.Add(due)
						timeout.occurred = done
						observe(next, err, done)
					}
				}
				timeout.Unlock()
			}
		}
		o(observer, subscribeOn, subscriber)
	})
	return observable
}

//jig:template Observable<Foo> Timeout
//jig:needs Observable Timeout

// Timeout mirrors the source ObservableFoo, but issues an error notification if
// a particular period of time elapses without any emitted items.
// Timeout schedules a task on the scheduler passed to it during subscription.
func (o ObservableFoo) Timeout(timeout time.Duration) ObservableFoo {
	return o.AsObservable().Timeout(timeout).AsObservableFoo()
}

//jig:template Timer
//jig:needs ObservableInt

// Timer creates an ObservableInt that emits a sequence of integers (starting
// at zero) after an initialDelay has passed. Subsequent values are emitted
// using  a schedule of intervals passed in. If only the initialDelay is
// given, Timer will emit only once.
func Timer(initialDelay time.Duration, intervals ...time.Duration) ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		runner := subscribeOn.ScheduleFutureRecursive(initialDelay, func(self func(time.Duration)) {
			if subscriber.Subscribed() {
				observe(i, nil, false)
				if subscriber.Subscribed() {
					if len(intervals) > 0 {
						self(intervals[i%len(intervals)])
					}
				}
				i++
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}

//jig:template Timestamp<Foo>

type TimestampFoo struct {
	Value     foo
	Timestamp time.Time
}

//jig:template Observable<Foo> Timestamp
//jig:needs Timestamp<Foo>, ObservableTimestamp<Foo>

// Timestamp attaches a timestamp to each item emitted by an observable
// indicating when it was emitted.
func (o ObservableFoo) Timestamp() ObservableTimestampFoo {
	observable := func(observe TimestampFooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			if subscriber.Subscribed() {
				if !done {
					observe(TimestampFoo{next, subscribeOn.Now()}, nil, false)
				} else {
					var zero TimestampFoo
					observe(zero, err, done)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template TimeInterval<Foo>

type TimeIntervalFoo struct {
	Value    foo
	Interval time.Duration
}

//jig:template Observable<Foo> TimeInterval
//jig:needs TimeInterval<Foo>, ObservableTimeInterval<Foo>

// TimeInterval intercepts the items from the source Observable and emits in
// their place a struct that indicates the amount of time that elapsed between
// pairs of emissions.
func (o ObservableFoo) TimeInterval() ObservableTimeIntervalFoo {
	observable := func(observe TimeIntervalFooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		begin := subscribeOn.Now()
		observer := func(next foo, err error, done bool) {
			if subscriber.Subscribed() {
				if !done {
					now := subscribeOn.Now()
					observe(TimeIntervalFoo{next, now.Sub(begin).Round(time.Millisecond)}, nil, false)
					begin = now
				} else {
					var zero TimeIntervalFoo
					observe(zero, err, done)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}
