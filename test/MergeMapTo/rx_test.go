package MergeMapTo

import _ "github.com/reactivego/rx"

func Example_basic() {
	Range(1, 2).MergeMapTo(Range(3, 2)).Println()
	// Output:
	// 3
	// 4
	// 3
	// 4
}
