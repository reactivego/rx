package rx

import "github.com/reactivego/scheduler"

//jig:template Scheduler

// Scheduler is used to schedule tasks to support subscribing and observing.
type Scheduler = scheduler.Scheduler

//jig:template MakeTrampolineScheduler
//jig:needs Scheduler

func MakeTrampolineScheduler() Scheduler {
	return scheduler.MakeTrampoline()
}

//jig:template GoroutineScheduler
//jig:needs Scheduler

func GoroutineScheduler() Scheduler {
	return scheduler.Goroutine
}

//jig:template Observable<Foo> ObserveOn

// ObserveOn specifies a dispatch function to use for delivering values to the observer.
func (o ObservableFoo) ObserveOn(dispatch func(task func())) ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			task := func() {
				observe(next, err, done)
			}
			dispatch(task)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> SubscribeOn

// SubscribeOn specifies the scheduler an ObservableFoo should use when it is
// subscribed to.
func (o ObservableFoo) SubscribeOn(subscribeOn Scheduler) ObservableFoo {
	observable := func(observe FooObserver, _ Scheduler, subscriber Subscriber) {
		if subscribeOn.IsConcurrent() {
			subscriber.OnWait(nil)
		} else {
			subscriber.OnWait(subscribeOn.Wait)
		}
		o(observe, subscribeOn, subscriber)
	}
	return observable
}
