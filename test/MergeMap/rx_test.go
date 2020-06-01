package MergeMap

import _ "github.com/reactivego/rx"

func Example_basic() {
	RangeInt(1, 2).MergeMapInt(func(n int) ObservableInt { return RangeInt(n, 2) }).Println()
	// Output:
	// 1
	// 2
	// 2
	// 3
}
