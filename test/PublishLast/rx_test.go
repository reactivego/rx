package PublishLast

import (
	"testing"
	"time"

	_ "github.com/reactivego/rx"
	. "github.com/reactivego/rx/test"
)

func Example_publishLast() {
	publish := From("a", "b", "c").PublishLast()
	go publish.Connect().Wait()
	publish.Println()
	publish.Println()
	// Output:
	// c
	// c
}

func TestPublishLast(e *testing.T) {
	const ms = time.Millisecond

	Describ(e, "operator", func(t T) {
		I(t, "should emit last notification of a simple source Observable", func(t T) {
			schedulers := []Scheduler{MakeTrampolineScheduler(), GoroutineScheduler()}
			for _, scheduler := range schedulers {
				published := FromInt(1, 2, 3, 4, 5).PublishLast()
				published.Subscribe(Expec(t).MakeIntObservation("actual"), scheduler)
				published.Connect(scheduler)
				scheduler.Wait()
				Expec(t).IntObservation("actual").ToBe(5)
			}
		})
		I(t, "should do nothing if connect is not called despite subscriptions", func(t T) {
			schedulers := []Scheduler{MakeTrampolineScheduler(), GoroutineScheduler()}
			for _, scheduler := range schedulers {
				published := From(1, 2, 3, 4, 5).PublishLast()
				published.Timeout(1*ms).Subscribe(Expec(t).MakeObservation("actual"), scheduler)
				scheduler.Wait()
				Expec(t).Observation("actual").ToBeError(TimeoutOccured)
			}
		})
		I(t, "should multicast the same values to multiple observers", func(t T) {
			schedulers := []Scheduler{MakeTrampolineScheduler(), GoroutineScheduler()}
			for _, scheduler := range schedulers {
				source := Concat(Just(1), Timer(5*ms), From(2, 3), Timer(50*ms), Just(4))
				published := source.PublishLast()
				Timer(1*ms).MergeMapTo(published.Observable).Subscribe(Expec(t).MakeObservation("1 ms"), scheduler)
				Timer(10*ms).MergeMapTo(published.Observable).Subscribe(Expec(t).MakeObservation("10 ms"), scheduler)
				Timer(100*ms).MergeMapTo(published.Observable).Subscribe(Expec(t).MakeObservation("100 ms"), scheduler)
				err := published.Connect(scheduler).Wait()
				scheduler.Wait()
				Asser(t).NoError(err)
				Expec(t).Observation("1 ms").ToBe(4)
				Expec(t).Observation("10 ms").ToBe(4)
				Expec(t).Observation("100 ms").ToBe(4)
			}
		})
		I(t, "should multicast an error from the source to multiple observers", func(t T) {
			schedulers := []Scheduler{MakeTrampolineScheduler(), GoroutineScheduler()}
			for _, scheduler := range schedulers {
				bazinga := RxError("bazinga")
				source := Concat(From(1, 2, 3), Timer(110*ms), Just(4), Throw(bazinga))
				published := source.PublishLast()
				Timer(1*ms).MergeMapTo(published.Observable).Subscribe(Expec(t).MakeObservation("1 ms"), scheduler)
				Timer(10*ms).MergeMapTo(published.Observable).Subscribe(Expec(t).MakeObservation("10 ms"), scheduler)
				Timer(100*ms).MergeMapTo(published.Observable).Subscribe(Expec(t).MakeObservation("100 ms"), scheduler)
				published.Connect(scheduler)
				scheduler.Wait()
				Expec(t).Observation("1 ms").ToBeError(bazinga)
				Expec(t).Observation("10 ms").ToBeError(bazinga)
				Expec(t).Observation("100 ms").ToBeError(bazinga)
			}
		})
		I(t, "should not cast any values to multiple observers, when source is unsubscribed explicitly and early", func(t T) {
			schedulers := []Scheduler{MakeTrampolineScheduler(), GoroutineScheduler()}
			for _, scheduler := range schedulers {
				source := Concat(From(1, 2, 3), Timer(110*ms), Just(4))
				published := source.PublishLast()
				Timer(1*ms).MergeMapTo(published.Observable).Timeout(120*ms).Subscribe(Expec(t).MakeObservation("1 ms"), scheduler)
				Timer(10*ms).MergeMapTo(published.Observable).Timeout(120*ms).Subscribe(Expec(t).MakeObservation("10 ms"), scheduler)
				Timer(100*ms).MergeMapTo(published.Observable).Timeout(120*ms).Subscribe(Expec(t).MakeObservation("100 ms"), scheduler)
				sub := published.Connect(scheduler)
				scheduler.ScheduleFuture(50*ms, func() {
					sub.Unsubscribe()
				})
				scheduler.Wait()
				Expec(t).Observation("1 ms").ToBeError(TimeoutOccured)
				Expec(t).Observation("10 ms").ToBeError(TimeoutOccured)
				Expec(t).Observation("100 ms").ToBeError(TimeoutOccured)
			}
		})
		I(t, "should not break unsubscription chains when result is unsubscribed explicitly", func(t T) {
			schedulers := []Scheduler{MakeTrampolineScheduler() /*, GoroutineScheduler()*/}
			for _, scheduler := range schedulers {
				source := Concat(From(1, 2, 3), Timer(110*ms), Just(4))
				published := source.MergeMap(func(x interface{}) Observable {
					return Just(x)
				}).PublishLast()
				Timer(1*ms).MergeMapTo(published.Observable).Timeout(120*ms).Subscribe(Expec(t).MakeObservation("1 ms"), scheduler)
				Timer(10*ms).MergeMapTo(published.Observable).Timeout(120*ms).Subscribe(Expec(t).MakeObservation("10 ms"), scheduler)
				Timer(100*ms).MergeMapTo(published.Observable).Timeout(120*ms).Subscribe(Expec(t).MakeObservation("100 ms"), scheduler)

				sub := published.Connect(scheduler)
				scheduler.ScheduleFuture(105*ms, func() {
					sub.Unsubscribe()
				})
				scheduler.Wait()
				Expec(t).Observation("1 ms").ToBeError(TimeoutOccured)
				Expec(t).Observation("10 ms").ToBeError(TimeoutOccured)
				Expec(t).Observation("100 ms").ToBeError(TimeoutOccured)
			}
		})
		Contex(t, "with refCount()", func(t T) {
			I(t, "should connect when first subscriber subscribes", func(t T) {
				schedulers := []Scheduler{MakeTrampolineScheduler(), GoroutineScheduler()}
				for _, scheduler := range schedulers {
					subscriptions := 0
					empty := Defer(func() Observable {
						subscriptions++
						return Empty()
					})
					source := Concat(empty, From(1, 2, 3), Timer(110*ms), Just(4))
					published := source.PublishLast()
					replayed := published.RefCount()
					Asser(t).Equal(subscriptions, 0)
					Timer(1*ms).MergeMapTo(replayed).Subscribe(Expec(t).MakeObservation("1 ms"), scheduler)
					scheduler.ScheduleFuture(5*ms, func() {
						Asser(t).Equal(subscriptions, 1)
					})
					Timer(10*ms).MergeMapTo(replayed).Subscribe(Expec(t).MakeObservation("10 ms"), scheduler)
					Timer(100*ms).MergeMapTo(replayed).Subscribe(Expec(t).MakeObservation("100 ms"), scheduler)
					scheduler.Wait()
					Expec(t).Observation("1 ms").ToBe(4)
					Expec(t).Observation("10 ms").ToBe(4)
					Expec(t).Observation("100 ms").ToBe(4)
				}
			})
			I(t, "should disconnect when last subscriber unsubscribes", func(t T) {
				schedulers := []Scheduler{MakeTrampolineScheduler() /*, GoroutineScheduler()*/}
				for _, scheduler := range schedulers {
					subscriptions := 0
					empty := Defer(func() Observable {
						subscriptions++
						return Empty()
					})
					source := Concat(empty, From(1, 2, 3), Timer(20*ms), Just(4))
					published := source.PublishLast()
					replayed := published.RefCount()
					sub1 := Timer(1*ms).MergeMapTo(replayed).Subscribe(Expec(t).MakeObservation("1 ms"), scheduler)
					sub2 := Timer(10*ms).MergeMapTo(replayed).Subscribe(Expec(t).MakeObservation("10 ms"), scheduler)
					Asser(t).Equal(subscriptions, 0)
					scheduler.ScheduleFuture(2*ms, func() {
						Asser(t).Equal(subscriptions, 1)
					})
					scheduler.ScheduleFuture(11*ms, sub1.Unsubscribe)
					var sub Subscription
					scheduler.ScheduleFuture(12*ms, func() {
						sub = published.Connect()
						Asser(t).Equal(sub.Subscribed(), true)
					})
					scheduler.ScheduleFuture(13*ms, sub2.Unsubscribe)
					scheduler.ScheduleFuture(16*ms, func() {
						Asser(t).Equal(sub.Subscribed(), false)
					})
					scheduler.Wait()
					Expec(t).Observation("1 ms").ToBeActive()
					Expec(t).Observation("10 ms").ToBeActive()
				}
			})
			I(t, "should be retryable unlike RxJS", func(t T) {
				schedulers := []Scheduler{MakeTrampolineScheduler(), GoroutineScheduler()}
				for _, scheduler := range schedulers {
					bazinga := RxError("bazinga")
					subscriptions := 0
					throw := Defer(func() Observable {
						subscriptions++
						if subscriptions < 2 {
							return Throw(bazinga)
						}
						return Empty()
					})
					source := Concat(From(1, 2, 3), Timer(100*ms), Just(4), throw)
					published := source.PublishLast().RefCount().Retry(3)
					Timer(1*ms).MergeMapTo(published).Subscribe(Expec(t).MakeObservation("1 ms"), scheduler)
					Timer(10*ms).MergeMapTo(published).Subscribe(Expec(t).MakeObservation("10 ms"), scheduler)
					Timer(50*ms).MergeMapTo(published).Subscribe(Expec(t).MakeObservation("50 ms"), scheduler)
					scheduler.Wait()
					Expec(t).Observation("1 ms").ToBe(4)
					Expec(t).Observation("10 ms").ToBe(4)
					Expec(t).Observation("50 ms").ToBe(4)
					Asser(t).Equal(subscriptions > 1, true)
				}
			})
		})
		I(t, "should multicast an empty source", func(t T) {
			schedulers := []Scheduler{MakeTrampolineScheduler(), GoroutineScheduler()}
			for _, scheduler := range schedulers {
				subscriptions := 0
				empty := Defer(func() Observable {
					subscriptions++
					return Empty()
				})
				published := empty.PublishLast()
				published.Subscribe(Expec(t).MakeObservation("actual"), scheduler)
				sub := published.Connect(scheduler)
				Asser(t).Equal(sub.Subscribed(), true)
				Asser(t).Equal(subscriptions, 1)
				scheduler.Wait()
				Asser(t).Equal(sub.Subscribed(), false)
				Expec(t).Observation("actual").ToComplete()
			}
		})
		I(t, "should multicast a never source", func(t T) {
			schedulers := []Scheduler{MakeTrampolineScheduler(), GoroutineScheduler()}
			for _, scheduler := range schedulers {
				subscribed := false
				never := Defer(func() Observable {
					subscribed = true
					return Never()
				})
				published := never.PublishLast()
				published.Timeout(1*ms).Subscribe(Expec(t).MakeObservation("actual"), scheduler)
				sub := published.Connect(scheduler)
				Asser(t).Equal(sub.Subscribed(), true)
				Asser(t).Equal(subscribed, true)
				scheduler.Wait()
				Asser(t).Equal(sub.Subscribed(), true)
				Expec(t).Observation("actual").ToBeError(TimeoutOccured)
			}
		})
		I(t, "should multicast a throw source", func(t T) {
			schedulers := []Scheduler{MakeTrampolineScheduler(), GoroutineScheduler()}
			for _, scheduler := range schedulers {
				bazinga := RxError("bazinga")
				subscriptions := 0
				throw := Defer(func() Observable {
					subscriptions++
					return Throw(bazinga)
				})
				published := throw.PublishLast()
				published.Subscribe(Expec(t).MakeObservation("actual"), scheduler)
				sub := published.Connect(scheduler)
				Asser(t).Equal(sub.Subscribed(), true)
				Asser(t).Equal(subscriptions, 1)
				scheduler.Wait()
				Asser(t).Equal(sub.Wait(), bazinga)
				Expec(t).Observation("actual").ToBeError(bazinga)
			}
		})
		I(t, "should multicast one observable to multiple observers", func(t T) {
			schedulers := []Scheduler{MakeTrampolineScheduler(), GoroutineScheduler()}
			for _, scheduler := range schedulers {
				subscriptions := 0
				source := Defer(func() Observable {
					subscriptions++
					return From(1, 2, 3, 4)
				})
				connectable := source.PublishLast()
				connectable.Subscribe(Expec(t).MakeObservation("actual1"), scheduler)
				connectable.Subscribe(Expec(t).MakeObservation("actual2"), scheduler)
				Expec(t).Observation("actual1").ToBeActive()
				Expec(t).Observation("actual2").ToBeActive()
				connectable.Connect(scheduler)
				scheduler.Wait()
				Expec(t).Observation("actual1").ToBe(4)
				Expec(t).Observation("actual2").ToBe(4)
				Asser(t).Equal(subscriptions, 1)
			}
		})
	})
}
