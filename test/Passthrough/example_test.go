package Passthrough

import "fmt"

func Example_passthrough() {
	Range(1, 3).Passthrough().SubscribeNext(func(next int) { fmt.Println(next) })
	// Output:
	// 1
	// 2
	// 3
}
