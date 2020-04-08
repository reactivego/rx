package Timeout

import (
	"fmt"
	"time"
)

//jig:no-doc

func Example_timeout() {
	const _0ms = 0
	const _250ms = 250 * time.Millisecond
	const _500ms = 500 * time.Millisecond

	scheduler := TrampolineScheduler()
	// scheduler := GoroutineScheduler()
	start := scheduler.Now()

	active := true
	source := MakeTimedInt(_0ms, func(Next func(int), Error func(error), Complete func()) time.Duration {
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

	if err := timed.Println(); err == ErrTimeout {
		fmt.Println(ErrTimeout.Error())
	}

	elapsed := scheduler.Now().Sub(start)
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
