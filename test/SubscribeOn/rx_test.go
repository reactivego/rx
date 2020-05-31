package SubscribeOn

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

// SubscribeOn selects the scheduler to use for running the subscription task.
func Example_goroutine() {
	concurrent := GoroutineScheduler()

	source := FromInt(1, 2, 3, 4, 5).SubscribeOn(concurrent)

	observe := func(next int, err error, done bool) {
		switch {
		case err != nil:
			fmt.Println(err)
		case done:
			fmt.Println("complete")
		default:
			fmt.Println(next)
		}
	}

	source.Subscribe(observe).Wait()

	//Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// complete
}
