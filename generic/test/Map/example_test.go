package Map

import (
	"fmt"
)

func Example_mapString() {
	err := FromInts(1, 2, 3, 4).MapString(func(i int) string { return fmt.Sprintf("%d!", i) }).Println()
	fmt.Println(err)
	
	//Output:
	// 1!
	// 2!
	// 3!
	// 4!
	// <nil>
}

func Example_mapVector() {
	// type Vector []int
	
	FromVector(Vector{1, 2, 3}, Vector{5, 6, 7}).MapInt(func(x Vector) int {
		sum := 0
		for _, v := range x {
			sum += v
		}
		return sum
	}).Println()

	//Output:
	// 6
	// 18
}
