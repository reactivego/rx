package RefCount

import (
	"fmt"
	"time"
)

// This example is a variant of the example in the book "Introduction to Rx"
// about using RefCount.
func Example_introToRx() {
	observable := Interval(27 * time.Millisecond)

	// Print when a value is published.
	observable = observable.Do(func(next int) { fmt.Printf("Publishing %d\n", next) })

	observable = observable.Publish().RefCount()

	// Make all subscriptions to observable asynchronous
	observable = observable.SubscribeOn(NewGoroutine())

	fmt.Println(">> Subscribing")
	subscription := observable.SubscribeNext(func(next int) { fmt.Printf("subscription : %d\n", next) })

	// The observable is hot for the next 100 milliseconds. It then will go
	// cold, unless another observer subscribes in that period.
	time.Sleep(100 * time.Millisecond)

	fmt.Println(">> Unsubscribing")
	subscription.Unsubscribe()
	fmt.Println(">> Finished")

	// Output:
	// >> Subscribing
	// Publishing 0
	// subscription : 0
	// Publishing 1
	// subscription : 1
	// Publishing 2
	// subscription : 2
	// >> Unsubscribing
	// >> Finished
}

// An example showing multiple subscriptions on a multicasting Publish
// Connectable who's Connect is controlled by a RefCount operator.
func Example_refCountMultipleSubscriptions() {
	channel := make(chan int, 30)
	source := FromChanInt(channel).Publish().RefCount().SubscribeOn(NewGoroutine())

	printNext := func(n int) { fmt.Println(n) }
	sub1 := source.SubscribeNext(printNext)
	sub2 := source.SubscribeNext(printNext)

	// 3 goroutines are now starting, 1 for publishing and 2 for subscribing.
	time.Sleep(100 * time.Millisecond)

	// goroutines should have started by now, feed the channel.
	channel <- 1
	channel <- 2
	channel <- 3

	// allow data in the channel to propagate.
	time.Sleep(100 * time.Millisecond)

	// cancel the first subscription.
	sub1.Unsubscribe()

	// more data for the second subscription
	channel <- 4

	// closing the channel will complete the observable and close the
	// subscription.
	close(channel)
	sub2.Wait()

	// Output:
	// 1
	// 1
	// 2
	// 2
	// 3
	// 3
	// 4
}
