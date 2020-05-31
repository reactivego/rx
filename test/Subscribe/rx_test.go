package Subscribe

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

// Shows how to correctly subscribe to an observable by calling the Subscribe
// method.
func Example_basic() {
	// Create an observable of int values
	observable := CreateInt(func(N NextInt, E Error, C Complete, X Canceled) {
		N(456)
		C()
		// Error will not be delivered by Subscribe, because when Subscribe got
		// the Complete it immediately canceled the subscription.
		E(RxError("error"))
	})

	// Observe function. Note that done is true for both an error as well as
	// for a normal completion.
	observe := func(next int, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}

	}

	// Subscribe and return a subscription.
	subscription := observable.Subscribe(observe)

	// Now wait for the observable to finish.
	subscription.Wait()

	// Because the observable completed, the subscription was canceled inside
	// the Subscribe method.
	if subscription.Subscribed() {
		fmt.Println("error: subscription should have been canceled")
	}
	// Output:
	// 456
	// complete
}

// Shows how to directly subscribe to an observable by calling it as a function.
// Note that this example incorrectly delivers the error after complete.
// So don't do this in production code, use the Subscribe method.
func Example_direct() {
	// Create an observable of int values
	observable := CreateInt(func(N NextInt, E Error, C Complete, X Canceled) {
		N(456)
		C()
		// Error will be delivered because we are not using the Subscribe method
		// to subscribe to the observable.
		E(RxError("error"))
	})

	// Observe function. Note that done is true for both an error as well as
	// for a normal completion.
	observe := func(next int, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}

	// Non standard way of subscribing to an observable. But it does illustrate
	// that an observable is just a function that can be called directly.
	scheduler := MakeTrampolineScheduler()
	subscriber := NewSubscriber()
	observable(observe, scheduler, subscriber)

	// We could just call scheduler.Wait(), but we hook the scheduler into the
	// subscriber like it is done in the Subscribe method implementation.
	subscriber.OnWait(scheduler.Wait)
	// Now wait for the observable to finish.
	subscriber.Wait()

	// Although the observable completed, the subscription is still active.
	if subscriber.Canceled() {
		fmt.Println("error: subscriber should still be subscribed")
	}
	// Output:
	// 456
	// complete
	// error
}
