package Connect

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_connect() {
	// Use Goroutine scheduler to run tasks concurrently.
	concurrent := GoroutineScheduler()

	// Source is an ObservableInt
	source := FromInt(1, 2)

	// Publish is an IntMulticaster
	publish := source.Publish()

	// First subscriber (concurrent)
	publish.SubscribeOn(concurrent).
		MapString(func(v int) string {
			return fmt.Sprintf("value is %d", v)
		}).
		Subscribe(func(next string, err error, done bool) {
			if !done {
				fmt.Println("sub1", next)
			}
		})

	// Second subscriber (concurrent)
	publish.SubscribeOn(concurrent).
		MapBool(func(v int) bool {
			return v > 1
		}).
		Subscribe(func(next bool, err error, done bool) {
			if !done {
				fmt.Println("sub2", next)
			}
		})

	// Third subscriber (concurrent)
	publish.SubscribeOn(concurrent).
		Subscribe(func(next int, err error, done bool) {
			if !done {
				fmt.Println("sub3", next)
			}
		})

	// Connect will cause the publisher to subscribe to the source on a
	// trampoline scheduler.
	subscription := publish.Connect()

	// Wait will actually cause the trampoline scheduler hooked into the
	// subscription to perform the connect to the source. And feed the source
	// into the Publish multicaster.
	subscription.Wait()

	// Wait for all subscriptions, running as concurrent tasks, to finish.
	concurrent.Wait()
	fmt.Println("tasks =", concurrent.Count())
	// Unordered output:
	// sub1 value is 1
	// sub2 false
	// sub3 1
	// sub2 true
	// sub3 2
	// sub1 value is 2
	// tasks = 0
}
