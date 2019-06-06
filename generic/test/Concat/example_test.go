package Concat

import (
	"fmt"
)

func Example_concat() {
	oa := FromInts(0, 1, 2, 3)
	ob := FromInts(4, 5)
	oc := FromInts(6)
	od := FromInts(7, 8, 9)
	ConcatInt(oa, ob, oc).Concat(od).SubscribeNext(func(next int) {
		fmt.Printf("%d,", next)
	})

	//Output:
	// 0,1,2,3,4,5,6,7,8,9,
}
