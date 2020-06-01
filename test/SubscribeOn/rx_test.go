package SubscribeOn

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

// SubscribeOn selects the scheduler to use for running the subscription task.
func Example_subscribeOn() {
	concurrent := GoroutineScheduler()

	source := FromInt(1, 2, 3, 4, 5).SubscribeOn(concurrent)

	observer := func(next int, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}

	source.Subscribe(observer).Wait()
	//Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// complete
}
