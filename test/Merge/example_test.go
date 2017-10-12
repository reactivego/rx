package Merge

import (
	"fmt"
	"time"
)

func Example_mergeInt() {
	o1 := FromInts(1, 2, 3).Debounce(100 * time.Millisecond)
	o2 := MergeInt(o1, EmptyInt()).DoOnComplete(func() {
		fmt.Println("a")
	})

	o3 := FromInts(4, 5, 6)
	o4 := MergeInt(o3, EmptyInt()).DoOnComplete(func() {
		fmt.Println("b")
	})

	MergeInt(o2, o4).DoOnComplete(func() {
		fmt.Println("c")
	}).Wait()

	//Output:
	// a
	// b
	// c
}
