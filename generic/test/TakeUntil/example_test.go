package TakeUntil

import "time"

func Example_takeUntil() {
	const _100ms = 100 * time.Millisecond
	const _250ms = 250 * time.Millisecond

	// emit "stop" after 250ms
	interrupt := Never().Timeout(_250ms).Catch(Just("stop")).SubscribeOn(GoroutineScheduler())

	// produce a number starting at 0 every 100ms until interrupted
	Interval(_100ms).TakeUntil(interrupt).Println()

	// Output:
	// 0
	// 1
}
