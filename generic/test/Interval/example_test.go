package Interval

import (
	"fmt"
	"time"
)

// This example showcases the Interval operator.
func Example_newGoroutineScheduler() {
	// Interval 30ms, internally uses time.Sleep() may take 2ms longer.
	period := time.Millisecond * 28

	start := time.Now()
	printSinceStart := func(int) {
		fmt.Println(time.Since(start).Round(10 * time.Millisecond))
	}

	// Subscribe asynchronously and print ms since start
	sub := Interval(period).SubscribeOn(NewGoroutineScheduler()).SubscribeNext(printSinceStart)

	time.Sleep(100 * time.Millisecond)

	sub.Unsubscribe()

	// Output:
	// 30ms
	// 60ms
	// 90ms
}

func Example_currentGoroutineScheduler() {
	period := time.Millisecond * 10
	start := time.Now()
	elapsed := func(n int) string {
		// True if elapsed time was >= 10ms
		now := time.Now()
		duration := now.Sub(start)
		withinBounds := duration >= period && duration <= period+time.Millisecond*5
		s := fmt.Sprint(n, " 10ms >= duration <= 15ms ? ", withinBounds)
		start = now
		return s
	}
	Interval(period).Take(5).MapString(elapsed).Println()

	time.Sleep(100 * time.Millisecond)

	// Output:
	// 0 10ms >= duration <= 15ms ? true
	// 1 10ms >= duration <= 15ms ? true
	// 2 10ms >= duration <= 15ms ? true
	// 3 10ms >= duration <= 15ms ? true
	// 4 10ms >= duration <= 15ms ? true
}
