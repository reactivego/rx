package Timeout

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_timeout() {
	const _0ms = 0
	const _250ms = 250 * time.Millisecond
	const _500ms = 500 * time.Millisecond

	scheduler := MakeTrampolineScheduler()
	//scheduler := GoroutineScheduler()
	start := scheduler.Now()

	active := true
	source := CreateFutureRecursiveInt(_0ms, func(Next NextInt, E Error, Complete Complete) time.Duration {
		if active {
			fmt.Println("Next(1) ->")
			Next(1)
			active = false
		} else {
			fmt.Println("Complete ->")
			Complete()
		}
		return _500ms
	})

	timed := source.Timeout(_250ms).SubscribeOn(scheduler)

	if err := timed.Println(); err == TimeoutOccured {
		fmt.Println(TimeoutOccured.Error())
	}

	elapsed := scheduler.Since(start)
	if elapsed > _250ms && elapsed < _500ms {
		fmt.Println("elapsed time is be between 250 and 500 ms")
	} else {
		fmt.Println("elapsed", elapsed)
	}

	// Output:
	// Next(1) ->
	// 1
	// timeout
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
	// timeout
	// 1
	// <nil>
}
