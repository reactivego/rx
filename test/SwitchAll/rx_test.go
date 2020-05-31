package SwitchAll

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_goroutine() {
	err := Interval(42 * time.Millisecond).
		Take(4).
		MapObservableInt(func(i int) ObservableInt {
			return Interval(16 * time.Millisecond).Take(4)
		}).
		SwitchAll().
		SubscribeOn(GoroutineScheduler()).
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
	err := Interval(42 * time.Millisecond).
		Take(4).
		MapObservableInt(func(i int) ObservableInt {
			return Interval(16 * time.Millisecond).Take(4)
		}).
		SwitchAll().
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
