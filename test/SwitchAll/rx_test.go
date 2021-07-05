package SwitchAll

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_goroutine() {
	type any = interface{}
	const ms = time.Millisecond
	err := Interval(42 * ms).
		Take(4).
		MapObservable(func(i any) Observable { return Interval(16 * ms).Take(4) }).
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
	const ms = time.Millisecond
	err := IntervalInt(42 * ms).
		Take(4).
		MapObservableInt(func(i int) ObservableInt { return IntervalInt(16 * ms).Take(4) }).
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
