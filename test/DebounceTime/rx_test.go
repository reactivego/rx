package DebounceTime

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx/generic"
)

// DebounceTime sees 100 emissions [0..99] at 1ms intervals.
// Then an interval of 11ms, so long enough for DebounceTime to emit 99.
func Example_debounceTime() {
	const ms = time.Millisecond

	interval := ConcatInt(IntervalInt(1*ms).Take(100), IntervalInt(11*ms).Take(1))

	interval.DebounceTime(10 * ms).Println()
	// Output: 99
}

// DebounceTime does not see a 10ms period in which nothing is emitted.
// So in this case it does not emit anything.
func Example_debounceTimeNoEmit() {
	const ms = time.Millisecond

	Interval(1 * ms).Take(100).DebounceTime(10 * ms).Println()
	// Output:
}

func Example_debounceTimeBursts() {
	const ms = time.Millisecond

	burst := func(i int) ObservableInt {
		return IntervalInt(5 * ms).Take(4).MapInt(func(j int) int {
			return i*100 + j
		})
	}

	fmt.Println("-1-")

	IntervalInt(100 * ms).Take(4).MergeMapInt(burst).DebounceTime(20 * ms).Println()

	fmt.Println("-2-")

	IntervalInt(20 * ms).Take(4).MergeMapInt(burst).DebounceTime(20 * ms).Println()
	// Output:
	// -1-
	// 3
	// 103
	// 203
	// -2-
}
