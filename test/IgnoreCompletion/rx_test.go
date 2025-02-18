package IgnoreCompletion

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx/generic"
)

func Example_ignoreComplete() {
	source := RangeInt(1, 5).IgnoreCompletion()

	// NOTE: subscription must run concurrently with main goroutine
	concurrent := GoroutineScheduler()
	subscription := source.SubscribeOn(concurrent).Subscribe(func(next int, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	})

	time.Sleep(100 * time.Millisecond)

	if subscription.Subscribed() {
		fmt.Println("subscription alive, as expected")
	}
	subscription.Unsubscribe()
	subscription.Wait()
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// subscription alive, as expected
}

func Example_ignoreError() {
	source := CreateInt(func(N NextInt, E Error, C Complete, X Canceled) {
		for i := 1; i < 6; i++ {
			N(i)
		}
		E(RxError("error"))
	}).IgnoreCompletion()

	// NOTE: subscription must run concurrently with main goroutine
	concurrent := GoroutineScheduler()
	subscription := source.SubscribeOn(concurrent).Subscribe(func(next int, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	})

	time.Sleep(100 * time.Millisecond)

	if subscription.Subscribed() {
		fmt.Println("subscription alive, as expected")
	}
	subscription.Unsubscribe()
	subscription.Wait()
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// subscription alive, as expected
}
