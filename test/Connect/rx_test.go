package Connect

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_connect() {
	// source is an ObservableInt
	source := FromInt(1, 2)

	// publish is an IntMulticaster
	publish := source.Publish()

	// concurrent will run a task concurrently on a new goroutine.
	concurrent := GoroutineScheduler()

	// First subscriber (concurrent)
	sub1 := publish.SubscribeOn(concurrent).
		MapString(func(v int) string {
			return fmt.Sprintf("value is %d", v)
		}).
		Subscribe(func(next string, err error, done bool) {
			if !done {
				fmt.Println("sub1", next)
			}
		})

	// Second subscriber (concurrent)
	sub2 := publish.SubscribeOn(concurrent).
		MapBool(func(v int) bool {
			return v > 1
		}).
		Subscribe(func(next bool, err error, done bool) {
			if !done {
				fmt.Println("sub2", next)
			}
		})

	// Third subscriber (concurrent)
	sub3 := publish.SubscribeOn(concurrent).
		Subscribe(func(next int, err error, done bool) {
			if !done {
				fmt.Println("sub3", next)
			}
		})

	// Connect will cause the publisher to subscribe to the source
	// Wait will actually cause the trampoline scheduler hooked into the
	// subscription to perform the connect to the source. And feed the
	// source into the Publish
	publish.Connect().Wait()

	// Wait for all subscriptions (running concurrently) to finish.
	sub1.Wait()
	sub2.Wait()
	sub3.Wait()

	// Unordered output:
	// sub1 value is 1
	// sub2 false
	// sub3 1
	// sub2 true
	// sub3 2
	// sub1 value is 2
}
