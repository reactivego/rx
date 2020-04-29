package Concat

import _ "github.com/reactivego/rx"

func Example_basic() {
	oa := FromInt(0, 1, 2, 3)
	ob := FromInt(4, 5)
	oc := FromInt(6)
	od := FromInt(7, 8, 9)
	ConcatInt(oa, ob, oc).Concat(od).Println()

	//Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
}
