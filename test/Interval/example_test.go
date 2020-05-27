package Interval

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_intervalAccuracy() {
	const ms = time.Millisecond
	start := time.Now()

	Interval(10 * ms).Take(1).Wait()
	
	onems := (10 * ms * ms) / time.Since(start)
	fmt.Println("1ms takes between 0.995ms and 1.005ms =", 995 < onems.Microseconds() && onems.Microseconds() < 1005)
	// Output: 1ms takes between 0.995ms and 1.005ms = true
}

// After 100ms the example cancels the observable subscription by calling the
// Unsubscribe method on the returned subscription.
func Example_intervalGoroutineScheduler() {
	const ms = time.Millisecond

	// Goroutine concurrently runs each subscribe on a separate goroutine.
	concurrent := GoroutineScheduler()

	// Print time since start of the subcribe.
	start := time.Now()
	printSinceStart := func(next int, err error, done bool) {
		if !done {
			elapsed := time.Since(start)
			fmt.Println(elapsed.Round(ms))
		}
	}

	// Interval 30ms. Subscribe and print ms since start.
	subscription := Interval(30 * ms).SubscribeOn(concurrent).Subscribe(printSinceStart)

	// Sleep for 100ms
	time.Sleep(100 * ms)

	// Cancel the subscription, Interval doesn't stop by itself.
	subscription.Unsubscribe()

	// Wait for the subscription (goroutine) to terminate.
	subscription.Wait()

	// Output:
	// 30ms
	// 60ms
	// 90ms
}

// Interval operator used together with the default Trampoline scheduler.
func Example_intervalTrampolineScheduler() {
	const ms = time.Millisecond
	const us = time.Microsecond

	start := time.Now()
	calculateElapsed := func(n int) string {
		// calculate elapsed time
		now := time.Now()
		elapsed := now.Sub(start)
		start = now

		// withinBounds true when elapsed time is between 9950us and 10050us
		withinBounds := (elapsed >= 9900 * us && elapsed <= 10100 * us)

		return fmt.Sprint(n, ": 9.9ms >= duration <= 10.1ms ? ", withinBounds)
	}

	Interval(10 * ms).Take(5).MapString(calculateElapsed).Println()

	// Output:
	// 0: 9.9ms >= duration <= 10.1ms ? true
	// 1: 9.9ms >= duration <= 10.1ms ? true
	// 2: 9.9ms >= duration <= 10.1ms ? true
	// 3: 9.9ms >= duration <= 10.1ms ? true
	// 4: 9.9ms >= duration <= 10.1ms ? true
}
