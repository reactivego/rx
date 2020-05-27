package Throttle

import (
	"time"

	_ "github.com/reactivego/rx"
)

func Example_throttle() {
	const ms = time.Millisecond

	Interval(1 * ms).Throttle(10 * ms).Take(5).Println()
	// Output:
	// 0
	// 10
	// 20
	// 30
	// 40
}
