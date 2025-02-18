package Timeout

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx/generic"
)

func Example_timeout() {
	const ms = time.Millisecond

	scheduler := NewScheduler()
	//scheduler := GoroutineScheduler()
	start := scheduler.Now()

	active := true
	source := CreateFutureRecursiveInt(0*ms, func(Next NextInt, E Error, Complete Complete) time.Duration {
		if active {
			fmt.Println("Next(1) ->")
			Next(1)
			active = false
		} else {
			fmt.Println("Complete ->")
			Complete()
		}
		return 500 * ms
	})

	timed := source.Timeout(250 * ms).SubscribeOn(scheduler)

	if err := timed.Println(); err == TimeoutOccured {
		fmt.Println(TimeoutOccured.Error())
	}

	elapsed := scheduler.Since(start)
	if elapsed > 250*ms && elapsed < 500*ms {
		fmt.Println("elapsed time is be between 250 and 500 ms")
	} else {
		fmt.Println("elapsed", elapsed)
	}

	// Output:
	// Next(1) ->
	// 1
	// timeout occured
	// elapsed time is be between 250 and 500 ms
}

func Example_timeoutTwice() {
	const ms = time.Millisecond

	subscription := 0
	source := DeferInt(func() ObservableInt {
		if subscription == 0 {
			subscription++
			return NeverInt()
		}
		return FromInt(1)
	})

	fmt.Println(source.Timeout(10 * ms).Println())
	fmt.Println(source.Timeout(10 * ms).Println())

	// Output:
	// timeout occured
	// 1
	// <nil>
}
