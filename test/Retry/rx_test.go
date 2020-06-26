package Retry

import (
	"fmt"
	"testing"

	_ "github.com/reactivego/rx"
	. "github.com/reactivego/rx/test"
)

func Example_retry() {
	errored := false
	a := CreateInt(func(N NextInt, E Error, C Complete, X Canceled) {
		N(1)
		N(2)
		N(3)
		if errored {
			C()
		} else {
			// Error triggers subscribe and subscribe is scheduled on trampoline....
			E(RxError("error"))
			errored = true
		}
	})
	err := a.Retry().Println()
	fmt.Println(errored)
	fmt.Println(err)
	// Output:
	// 1
	// 2
	// 3
	// 1
	// 2
	// 3
	// true
	// <nil>
}

func Example_retryConcurrent() {
	scheduler := GoroutineScheduler()
	errored := false
	a := CreateInt(func(N NextInt, E Error, C Complete, X Canceled) {
		N(1)
		N(2)
		N(3)
		if errored {
			C()
		} else {
			E(RxError("error"))
			errored = true
		}
	})
	err := a.Retry().SubscribeOn(scheduler).Println()
	fmt.Println(errored)
	fmt.Println(err)
	// Output:
	// 1
	// 2
	// 3
	// 1
	// 2
	// 3
	// true
	// <nil>
}

func TestRetryPublishLatest(e *testing.T) {
	Describ(e, "Retry", func(t T) {
		I(t, "can retry a published source", func(t T) {
			schedulers := []Scheduler{MakeTrampolineScheduler(), GoroutineScheduler()}
			for _, scheduler := range schedulers {

				subscriptions := 0
				source := Defer(func() Observable {
					subscriptions++
					if subscriptions == 1 {
						return Throw(RxError("error"))
					}
					return From(1, 2, 3)
				})
				observable := source.PublishLast().RefCount().Retry(3)
				expected := []interface{}{3}

				Expect := Expec(t)
				observable.Subscribe(Expect.MakeObservation("1"), scheduler)
				observable.Subscribe(Expect.MakeObservation("2"), scheduler)
				observable.Subscribe(Expect.MakeObservation("3"), scheduler)
				observable.Subscribe(Expect.MakeObservation("4"), scheduler)
				scheduler.Wait()
				// Asser(t).Equal(subscriptions, 2)
				Expect.Observation("1").ToBe(expected)
				Expect.Observation("2").ToBe(expected)
				Expect.Observation("3").ToBe(expected)
				Expect.Observation("4").ToBe(expected)
			}
		})
	})
}
