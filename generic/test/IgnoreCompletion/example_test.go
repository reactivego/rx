package IgnoreCompletion

import (
	"fmt"
	"time"
)

func Example_ignoreCompletion() {
	source := Range(1, 5).IgnoreCompletion()

	// Subscribe asynchronously
	subscription := source.SubscribeNext(func(next int) {
		fmt.Println(next)
	}, SubscribeOn( /*must be async*/ GoroutineScheduler()))

	time.Sleep(100 * time.Millisecond)

	if !subscription.Closed() {
		fmt.Println("subscription alive, as expected")
	}
	subscription.Unsubscribe()

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// subscription alive, as expected
}

func Example_ignoreCompletionError() {
	source := CreateInt(func(observer IntObserver) {
		for i := 1; i < 6; i++ {
			observer.Next(i)
		}
		observer.Error(RxError("error"))
	}).IgnoreCompletion()

	// Subscribe asynchronously
	subscription := source.SubscribeNext(func(next int) {
		fmt.Println(next)
	}, SubscribeOn( /*must be async*/ GoroutineScheduler()))

	time.Sleep(100 * time.Millisecond)

	if !subscription.Closed() {
		fmt.Println("subscription alive, as expected")
	}
	subscription.Unsubscribe()

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// subscription alive, as expected
}
