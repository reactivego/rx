package RefCount

import (
	"fmt"
	"os"
	"time"
	"sync"
)

// This example is a variant of the example in the book "Introduction to Rx"
// about using RefCount.
func Example_introToRx() {
	// Subscribe on concurrent scheduler
	scheduler := NewGoroutineScheduler()

	observable := Interval(50 * time.Millisecond)

	// Print when a value is published.
	observable = observable.Do(func(next int) { fmt.Printf("published: %d\n", next) })

	observable = observable.Publish().RefCount()

	// Make all subscriptions to observable concurrent
	observable = observable.SubscribeOn(scheduler)

	fmt.Println(">> Subscribing")
	subscription := observable.SubscribeNext(func(next int) { fmt.Printf("subscription : %d\n", next) })

	// The observable is hot for the next 100 milliseconds. It then will go
	// cold, unless another observer subscribes in that period.
	time.Sleep(175 * time.Millisecond)

	fmt.Println(">> Unsubscribing")
	subscription.Unsubscribe()
	fmt.Println(">> Finished")

	// Output:
	// >> Subscribing
	// published: 0
	// subscription : 0
	// published: 1
	// subscription : 1
	// published: 2
	// subscription : 2
	// >> Unsubscribing
	// >> Finished
}

// An example showing multiple subscriptions on a multicasting Publish
// Connectable who's Connect is controlled by a RefCount operator.
func Example_refCountMultipleSubscriptions() {
	// Subscribe on concurrent scheduler
	scheduler := NewGoroutineScheduler()

	var wg sync.WaitGroup
	channel := make(chan int, 30)
	source := FromChanInt(channel).Publish().RefCount().SubscribeOn(scheduler)

	sub1 := source.SubscribeNext(func(n int) {
		fmt.Println(n)
		os.Stdout.Sync()
		wg.Done()
	})
	sub2 := source.SubscribeNext(func(n int) {
		fmt.Println(n)
		os.Stdout.Sync()
		wg.Done()
	})

	// 3 goroutines are now starting, 1 for publishing and 2 for subscribing.
	// in the mean time we start feeding the channel.
	wg.Add(6)
	channel <- 1
	channel <- 2
	channel <- 3

	// wait for the channel data to propagate
	wg.Wait()

	// cancel the first subscription.
	sub1.Unsubscribe()

	// more data for the second subscription, then close the
	// channel to complete the observable and make Wait return.
	wg.Add(1)
	channel <- 4
	close(channel)

	sub2.Wait()
	wg.Wait()

	// Output:
	// 1
	// 1
	// 2
	// 2
	// 3
	// 3
	// 4
}
