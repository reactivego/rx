package Audit

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_audit() {
	const ms = time.Millisecond

	Interval(1 * ms).Audit(10 * ms).Take(5).Println()
	// Output:
	// 9
	// 19
	// 29
	// 39
	// 49
}

func Example_auditBursts() {
	const ms = time.Millisecond

	burst20ms := func(i int) ObservableInt {
		return Interval(5 * ms).Take(4).MapInt(func(j int) int {
			return i*100 + j
		})
	}

	fmt.Println("-1-")

	Interval(25 * ms).Take(4).MergeMapInt(burst20ms).Audit(21 * ms).Println()

	fmt.Println("-2-")

	Interval(50 * ms).Take(4).MergeMapInt(burst20ms).Audit(14 * ms).Println()
	// Output:
	// -1-
	// 3
	// 103
	// 203
	// -2-
	// 2
	// 3
	// 102
	// 103
	// 202
	// 203
	// 302
}
