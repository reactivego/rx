package Concat

import _ "github.com/reactivego/rx"

func Example_basic() {
	oa := FromInt(0, 1, 2, 3)
	ob := FromInt(4, 5)
	oc := FromInt(6)
	ConcatInt(oa, ob, oc).Println()

	//Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
}
