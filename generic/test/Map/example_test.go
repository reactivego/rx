package Map

import (
	"fmt"
)

func Example_mapString() {
	result, err := FromInts(1, 2, 3, 4).MapString(func(i int) string { return fmt.Sprintf("%d!", i) }).ToSlice()
	if err != nil {
		panic(err)
	}
	for _, v := range result {
		fmt.Println(v)
	}

	//Output:
	// 1!
	// 2!
	// 3!
	// 4!
}

type Vector []int

func Example_mapVector() {

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
