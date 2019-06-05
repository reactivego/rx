package rx

import (
	"github.com/reactivego/scheduler"
)

//jig:template Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler interface {
	Schedule(task func())
}

//jig:template NewScheduler

func NewGoroutineScheduler() Scheduler  { return &scheduler.Goroutine{} }
func NewTrampolineScheduler() Scheduler { return &scheduler.Trampoline{} }

//jig:template Observable<Foo> ObserveOn
//jig:needs NewScheduler

// ObserveOn specifies the scheduler on which an observer will observe this
// ObservableFoo.
func (o ObservableFoo) ObserveOn(observeOn Scheduler) ObservableFoo {
	observable := func(observe FooObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			observeOn.Schedule(func() {
				observe(next, err, done)
			})
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> SubscribeOn
//jig:needs NewScheduler

// SubscribeOn specifies the scheduler an ObservableFoo should use when it is
// subscribed to.
func (o ObservableFoo) SubscribeOn(subscribeOn Scheduler) ObservableFoo {
	observable := func(observe FooObserveFunc, _ Scheduler, subscriber Subscriber) {
		o(observe, subscribeOn, subscriber)
	}
	return observable
}
