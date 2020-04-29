package MergeMap

import _ "github.com/reactivego/rx"

func Example_basic() {
	Range(1, 2).MergeMapInt(func(n int) ObservableInt { return Range(n, 2) }).Println()
	// Output:
	// 1
	// 2
	// 2
	// 3
}
