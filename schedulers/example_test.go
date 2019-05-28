package schedulers

import (
	"fmt"
)

func ExampleTrampoline() {
	tramp := &Trampoline{}
	fmt.Println("before")
	tramp.Dispatch(func() {
		fmt.Println("1")

		tramp.Dispatch(func() {
			fmt.Println("3")

			tramp.Dispatch(func() {
				fmt.Println("5")
			})

			fmt.Println("4")
		})

		fmt.Println("2")
	})
	fmt.Println("after")

	// Output:
	// before
	// 1
	// 2
	// 3
	// 4
	// 5
	// after
}
