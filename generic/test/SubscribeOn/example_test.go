package SubscribeOn

import (
	"fmt"
)

// SubscribeOn can be used to change the scheduler the Subcribe method uses to
// subscribe to the source observable.
func Example_goroutine() {
	scheduler := NewGoroutineScheduler()
	source := FromInts(1, 2, 3, 4, 5).SubscribeOn(scheduler)

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

	// scheduler is asynchronous, so Subscribe is as well.
	subscription := source.Subscribe(observe)
	fmt.Println("subscribed", !subscription.Closed())
	// now need to wait for the subscription to terminate.
	subscription.Wait()

	//Output:
	// subscribed true
	// 1
	// 2
	// 3
	// 4
	// 5
	// complete
}
