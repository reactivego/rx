package rx

import (
	"fmt"
)

func ExampleObservable_Concat() {
	oa := From(0, 1, 2, 3)
	ob := From(4, 5)
	oc := From(6)
	od := From(7, 8, 9)
	oa.Concat(ob, oc).Concat(od).SubscribeNext(func(next any) {
		fmt.Printf("%d,", next.(int))
	})

	// Output:
	// 0,1,2,3,4,5,6,7,8,9,
}
