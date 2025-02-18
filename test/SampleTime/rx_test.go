package SampleTime

import (
	"time"

	_ "github.com/reactivego/rx/generic"
)

func Example_sampleTime() {
	const ms = time.Millisecond

	Interval(90 * ms).SampleTime(200 * ms).Take(3).Println()
	// Output:
	// 1
	// 3
	// 5
}
