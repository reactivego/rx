package Interval

import (
	"fmt"
	"time"
)

// This example showcases the Interval operator. After 100ms the example
// cancels the observable subscription by calling the Unsubscribe method on the
// returned subscription.
func Example_newGoroutineScheduler() {
	const _10ms = 10 * time.Millisecond
	const _28ms = 28 * time.Millisecond

	// NewGoroutine asynchronously runs each subscribe on a separate goroutine.
	scheduler := NewGoroutineScheduler()

	// Print time since start of the subcribe.
	start := time.Now()
	printSinceStart := func(int) {
		elapsed := time.Since(start)
		fmt.Println(elapsed.Round(_10ms))
	}

	// Interval 30ms usually takes up to 4ms longer, so we'er using 28ms.
	// Subscribe (asynchronous through NewGoroutine) and print ms since start.
	subscription := Interval(_28ms).SubscribeOn(scheduler).SubscribeNext(printSinceStart)

	// Sleep for 100ms
	time.Sleep(100 * time.Millisecond)

	// Cancel the subscription, Interval doesn't stop by itself.
	subscription.Unsubscribe()

	// Output:
	// 30ms
	// 60ms
	// 90ms
}

// Interval operator used together with the default CurrentGoroutine scheduler.
// 
func Example_currentGoroutineScheduler() {
	const _10ms = 10 * time.Millisecond
	const _15ms = 15 * time.Millisecond

	start := time.Now()
	calculateElapsed := func(n int) string {
		// calculate elapsed time
		now := time.Now()
		elapsed := now.Sub(start)
		start = now

		// withinBounds true when elapsed time is between 10ms and 15ms
		withinBounds := (elapsed >= _10ms && elapsed <= _15ms)

		return fmt.Sprint(n, ": 10ms >= duration <= 15ms ? ", withinBounds)
	}

	Interval(_10ms).Take(5).MapString(calculateElapsed).Println()

	// Output:
	// 0: 10ms >= duration <= 15ms ? true
	// 1: 10ms >= duration <= 15ms ? true
	// 2: 10ms >= duration <= 15ms ? true
	// 3: 10ms >= duration <= 15ms ? true
	// 4: 10ms >= duration <= 15ms ? true
}
