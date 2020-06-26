package rx

import "time"

//jig:template Observable<Foo> Catch

// Catch recovers from an error notification by continuing the sequence without
// emitting the error but by switching to the catch ObservableFoo to provide
// items.
func (o ObservableFoo) Catch(catch ObservableFoo) ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			if err != nil {
				catch(observe, subscribeOn, subscriber)
			} else {
				observe(next, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> CatchError

// CatchError catches errors on the ObservableFoo to be handled by returning a
// new ObservableFoo or throwing an error. It is passed a selector function
// that takes as arguments err, which is the error, and caught, which is the
// source observable, in case you'd like to "retry" that observable by
// returning it again. Whatever observable is returned by the selector will be
// used to continue the observable chain.
func (o ObservableFoo) CatchError(selector func(err error, caught ObservableFoo) ObservableFoo) ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			if err != nil {
				selector(err, o)(observe, subscribeOn, subscriber)
			} else {
				observe(next, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Retry

// Retry if a source ObservableFoo sends an error notification, resubscribe to
// it in the hopes that it will complete without error. If count is zero or
// negative, the retry count will be effectively infinite. The scheduler
// passed when subscribing is used by Retry to schedule any retry attempt. The
// time between retries is 1 millisecond, so retry frequency is 1 kHz. Any
// SubscribeOn operators should be called after Retry to prevent lockups
// caused by mixing different schedulers in the same subscription for retrying
// and subscribing.
func (o ObservableFoo) Retry(count ...int) ObservableFoo {
	count = append(count, int(^uint(0)>>1))
	if count[0] <= 0 {
		count = []int{int(^uint(0) >> 1)}
	}
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var retry struct {
			count       int
			observer    FooObserver
			subscriber  Subscriber
			resubscribe func()
		}
		retry.count = count[0]
		retry.observer = func(next foo, err error, done bool) {
			if err != nil && retry.count > 0 {
				retry.count--
				retry.subscriber.Done(err)
				subscribeOn.ScheduleFuture(1*time.Millisecond, retry.resubscribe)
			} else {
				observe(next, err, done)
			}
		}
		retry.resubscribe = func() {
			if subscriber.Subscribed() {
				retry.subscriber = subscriber.Add()
				if !subscribeOn.IsConcurrent() {
					retry.subscriber.OnWait(subscribeOn.Wait)
					subscriber.OnWait(func() { retry.subscriber.Wait() })
				}
				o(retry.observer, subscribeOn, retry.subscriber)
			}
		}
		retry.resubscribe()
	}
	return observable
}
