package PublishBehavior

import (
	"testing"
	"time"

	_ "github.com/reactivego/rx/generic"
	. "github.com/reactivego/rx/test"
)

func Example_publishBehavior() {
	publish := From("b", "c").PublishBehavior("a")
	go func() {
		time.Sleep(time.Microsecond)
		publish.Connect().Wait()
	}()
	publish.Println()
	publish.Println()
	publish.Println()
	publish.Println()
	// Output:
	// a
	// b
	// c
}

func TestPublishBehavior(e *testing.T) {
	const ms = time.Millisecond
	Describ(e, "operator", func(t T) {
		I(t, "should mirror a simple source Observable", func(t T) {
			schedulers := []Scheduler{GoroutineScheduler(), NewScheduler()}
			for _, scheduler := range schedulers {
				source := Concat(Timer(1*ms).MergeMapTo(From(1, 2)), Timer(1*ms).MergeMapTo(From(3, 4)), Timer(1*ms).MergeMapTo(Just(5)))
				published := source.PublishBehavior(0)
				published.Subscribe(Expec(t).MakeObservation("actual"), scheduler)
				published.Connect(scheduler)
				scheduler.Wait()
				Expec(t).Observation("actual").ToBe(0, 1, 2, 3, 4, 5)
				Expec(t).Observation("actual").ToComplete()
			}
		})
		I(t, "should only emit default value if connect is not called, despite subscription", func(t T) {
			schedulers := []Scheduler{GoroutineScheduler(), NewScheduler()}
			for _, scheduler := range schedulers {
				source := Concat(Timer(1*ms).MergeMapTo(From(1, 2)), Timer(1*ms).MergeMapTo(From(3, 4)), Timer(1*ms).MergeMapTo(Just(5)))
				published := source.PublishBehavior(0)
				published.Timeout(1*ms).Subscribe(Expec(t).MakeObservation("actual"), scheduler)
				scheduler.Wait()
				Expec(t).Observation("actual").ToBe(0)
				Expec(t).Observation("actual").ToBeError(TimeoutOccured)
			}
		})
		I(t, "should multicast the same values to multiple observers", func(t T) {
			schedulers := []Scheduler{GoroutineScheduler(), NewScheduler()}
			for _, scheduler := range schedulers {
				source := Concat(Timer(1*ms).MergeMapTo(From(1, 2)), Timer(2*ms).MergeMapTo(Just(3)), Timer(2*ms).MergeMapTo(Just(4)))
				published := source.PublishBehavior(0)
				Timer(0*ms).MergeMapTo(Just("a")).MergeMapTo(published.Observable).Subscribe(Expec(t).MakeObservation("0 ms"), scheduler)
				Timer(2*ms).MergeMapTo(Just("b")).MergeMapTo(published.Observable).Subscribe(Expec(t).MakeObservation("2 ms"), scheduler)
				Timer(4*ms).MergeMapTo(Just("b")).MergeMapTo(published.Observable).Subscribe(Expec(t).MakeObservation("4 ms"), scheduler)
				published.Connect(scheduler)
				scheduler.Wait()
				Expec(t).Observation("0 ms").ToBe(0, 1, 2, 3, 4)
				Expec(t).Observation("0 ms").ToComplete()
				Expec(t).Observation("2 ms").ToBe(2, 3, 4)
				Expec(t).Observation("2 ms").ToComplete()
				Expec(t).Observation("4 ms").ToBe(3, 4)
				Expec(t).Observation("4 ms").ToComplete()
			}
		})
		I(t, "should multicast an error from the source to multiple observers", func(t T) {
			schedulers := []Scheduler{GoroutineScheduler(), NewScheduler()}
			for _, scheduler := range schedulers {
				bazinga := RxError("bazinga")
				source := Concat(Timer(1*ms).MergeMapTo(From(1, 2)), Timer(2*ms).MergeMapTo(Just(3)), Timer(2*ms).MergeMapTo(Just(4)), Throw(bazinga))
				published := source.PublishBehavior(0)
				Timer(0*ms).MergeMapTo(Just("a")).MergeMapTo(published.Observable).Subscribe(Expec(t).MakeObservation("0 ms"), scheduler)
				Timer(2*ms).MergeMapTo(Just("b")).MergeMapTo(published.Observable).Subscribe(Expec(t).MakeObservation("2 ms"), scheduler)
				Timer(4*ms).MergeMapTo(Just("b")).MergeMapTo(published.Observable).Subscribe(Expec(t).MakeObservation("4 ms"), scheduler)
				published.Connect(scheduler)
				scheduler.Wait()
				Expec(t).Observation("0 ms").ToBe(0, 1, 2, 3, 4)
				Expec(t).Observation("0 ms").ToBeError(bazinga)
				Expec(t).Observation("2 ms").ToBe(2, 3, 4)
				Expec(t).Observation("2 ms").ToBeError(bazinga)
				Expec(t).Observation("4 ms").ToBe(3, 4)
				Expec(t).Observation("4 ms").ToBeError(bazinga)
			}
		})
		I(t, "should not break unsubscription chains when result is unsubscribed explicitly", func(t T) {})
		Contex(t, "with refCount()", func(t T) {
			I(t, "should connect when first subscriber subscribes", func(t T) {})
			I(t, "should disconnect when last subscriber unsubscribes", func(t T) {})
			I(t, "should NOT be retryable", func(t T) {})
			I(t, "should NOT be repeatable", func(t T) {})
			I(t, "should emit completed when subscribed after completed", func(t T) {})
			I(t, "should multicast an empty source", func(t T) {})
			I(t, "should multicast a never source", func(t T) {})
			I(t, "should multicast a throw source", func(t T) {})
			I(t, "should multicast one observable to multiple observers", func(t T) {})
			I(t, "should follow the RxJS 4 behavior and emit nothing to observer after completed", func(t T) {})
		})
	})
}
