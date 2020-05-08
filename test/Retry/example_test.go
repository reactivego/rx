package Retry

import (
	"fmt"

	_ "github.com/reactivego/rx"
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
	}).SubscribeOn(TrampolineScheduler())
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
	}).SubscribeOn(GoroutineScheduler())
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
