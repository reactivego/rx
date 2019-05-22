package Range

import (
	"fmt"
)

func Example_subscribe() {
	Range(1, 5).Subscribe(func(next int, err error, done bool) {
		if !done {
			fmt.Println(next)
		} else {
			fmt.Println("done")
		}
	})

	//Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// done
}
