package SubscribeOn

import (
	"fmt"
)

func Example_goroutine() {
	FromInts(1, 2, 3, 4, 5).SubscribeOn(NewGoroutine()).Subscribe(func(next int, err error, done bool) {
		switch {
		case err != nil:
			fmt.Println(err)
		case done:
			fmt.Println("complete")
		default:
			fmt.Println(next)
		}
	}).Wait()

	//Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// complete
}
