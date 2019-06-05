package rx

import (
	"time"
)


func ExampleObservableObservable_SwitchAll() {

	// SwitchAll does not work well with the default trampoline scheduler, so we use a goroutine scheduler instead.
	scheduler := NewGoroutineScheduler()

	// intToObs creates a new observable that emits an integer starting after and then repeated every 20 milliseconds
	// in the range starting at 0 and incrementing by 1. It takes only the first 10 emitted values and then uses
	// AsObservable to convert the IntObservable back to an untyped Observable.
	intToObs := func(i int) Observable {
		return Interval(20 * time.Millisecond).
				Take(10).
				AsObservable()
	}

	Interval(100 * time.Millisecond).
		Take(3).
		MapObservable(intToObs).
		SwitchAll().
		Println(SubscribeOn(scheduler))

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 0
	// 1
	// 2
	// 3
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
}

