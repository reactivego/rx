package ConcatWith

import _ "github.com/reactivego/rx"

func Example_concatWith() {
	oa := FromInt(0, 1, 2)
	ob := FromInt(3)
	oc := FromInt(4, 5)
	oa.ConcatWith(ob, oc).Println()

	//Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
}
