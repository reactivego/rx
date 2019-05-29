package TakeUntil

import (
	"time"
)

func Example_takeUntil() {
	scheduler := NewGoroutine()
	interrupt := Never().Timeout(150 * time.Millisecond).Catch(Just("stop")).SubscribeOn(scheduler)

	Interval(100 * time.Millisecond).TakeUntil(interrupt).Println()

	// Output:
	// 0
}
