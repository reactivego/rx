package Debounce

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)


// Debounce sees 100 emissions [0..99] at 1ms intervals.
// Then an interval of 11ms, so long enough for Debounce to emit 99.
func Example_debounce() {
	const ms = time.Millisecond

	interval := ConcatInt(Interval(1*ms).Take(100), Interval(11*ms).Take(1))

	interval.Debounce(10 * ms).Println()
	// Output: 99
}

// Debounce does not see a 10ms period in which nothing is emitted.
// So in this case it does not emit anything.
func Example_debounceNoEmit() {
	const ms = time.Millisecond

	Interval(1*ms).Take(100).Debounce(10 * ms).Println()
	// Output:
}

func Example_debounceBursts() {
	const ms = time.Millisecond

	burst := func(i int) ObservableInt {
		return Interval(5 * ms).Take(4).MapInt(func(j int) int {
			return i*100 + j
		})
	}

	fmt.Println("-1-")

	Interval(100 * ms).Take(4).MergeMapInt(burst).Debounce(20 * ms).Println()

	fmt.Println("-2-")

	Interval(20 * ms).Take(4).MergeMapInt(burst).Debounce(20 * ms).Println()
	// Output:
	// -1-
	// 3
	// 103
	// 203
	// -2-
}
