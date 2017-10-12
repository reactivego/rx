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
