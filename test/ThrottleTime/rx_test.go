package ThrottleTime

import (
	"time"

	_ "github.com/reactivego/rx/generic"
)

func Example_throttleTime() {
	const ms = time.Millisecond

	IntervalInt(1 * ms).ThrottleTime(10 * ms).Take(5).Println()
	// Output:
	// 0
	// 10
	// 20
	// 30
	// 40
}
