package Sample

import (
	"time"

	_ "github.com/reactivego/rx"
)

func Example_sample() {
	const ms = time.Millisecond

	Interval(90 * ms).Sample(200 * ms).Take(3).Println()
	// Output:
	// 1
	// 3
	// 5
}
