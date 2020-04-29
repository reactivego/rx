package Retry

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_retry() {
	errored := false
	a := CreateInt(func(observer IntObserver) {
		observer.Next(1)
		observer.Next(2)
		observer.Next(3)
		if errored {
			observer.Complete()
		} else {
			// Error triggers subscribe and subscribe is scheduled on trampoline....
			observer.Error(RxError("error"))
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
	a := CreateInt(func(observer IntObserver) {
		observer.Next(1)
		observer.Next(2)
		observer.Next(3)
		if errored {
			observer.Complete()
		} else {
			observer.Error(RxError("error"))
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
