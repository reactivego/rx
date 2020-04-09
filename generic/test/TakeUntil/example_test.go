package TakeUntil

import (
	"time"
)

func Example_takeUntil() {
	interrupt := Never().Timeout(150 * time.Millisecond).Catch(Just("stop")).SubscribeOn(GoroutineScheduler())

	Interval(100 * time.Millisecond).TakeUntil(interrupt).Println()

	// Output:
	// 0
}
