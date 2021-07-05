package TakeUntil

import (
	"time"

	_ "github.com/reactivego/rx"
)

func Example_takeUntil() {
	const ms = time.Millisecond

	// emit "stop" after 250ms
	interrupt := Never().Timeout(250 * ms).Catch(Just("stop")).SubscribeOn(GoroutineScheduler())

	// produce a number starting at 0 every 100ms until interrupted
	IntervalInt(100 * ms).TakeUntil(interrupt).Println()

	// Output:
	// 0
	// 1
}
