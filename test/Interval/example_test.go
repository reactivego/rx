package Interval

import (
	"fmt"
	"time"
)

// This example showcases the Interval operator.
func Example_interval() {
	start := time.Now()

	// Interval 30ms, internally uses time.Sleep() which on Mac usually takes
	// between 2ms to 3ms longer.
	interval := Interval(27 * time.Millisecond)

	printSinceStart := func(int) {
		fmt.Println(time.Since(start).Round(10 * time.Millisecond))
	}

	// Subscribe asynchronously and print ms since start
	sub := interval.SubscribeNext(printSinceStart, SubscribeOn(NewGoroutine()))

	time.Sleep(100 * time.Millisecond)
	sub.Unsubscribe()

	// Output:
	// 30ms
	// 60ms
	// 90ms
}
