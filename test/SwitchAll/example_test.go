package SwitchAll

import (
	"fmt"
	"time"
)

func Example_goroutine() {
	scheduler := GoroutineScheduler()
	err := Interval(42 * time.Millisecond).
		Take(4).
		MapObservableInt(func(i int) ObservableInt {
			return Interval(16 * time.Millisecond).Take(4)
		}).
		SwitchAll().
		SubscribeOn(scheduler).
		Println()

	if err == nil {
		fmt.Println("success")
	}

	// Output:
	// 0
	// 1
	// 0
	// 1
	// 0
	// 1
	// 0
	// 1
	// 2
	// 3
	// success
}

func Example_trampoline() {
	scheduler := TrampolineScheduler()
	err := Interval(42 * time.Millisecond).
		Take(4).
		MapObservableInt(func(i int) ObservableInt {
			return Interval(16 * time.Millisecond).Take(4)
		}).
		SwitchAll().
		SubscribeOn(scheduler).
		Println()

	if err == nil {
		fmt.Println("success")
	}

	// Output:
	// 0
	// 1
	// 0
	// 1
	// 0
	// 1
	// 0
	// 1
	// 2
	// 3
	// success
}
