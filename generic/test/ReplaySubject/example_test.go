package ReplaySubject

import (
	"fmt"
	"time"
)

// ReplaySubject is both an observer and an observable. The maximum buffersize
// is passed in along with how long entries in the replay buffer stay fresh.
// As long as entries are fresh they are send to an observer that subscribes.
func Example_replaySubject() {
	subject := NewReplaySubjectInt(10, time.Hour)

	// Feed the subject on a separate goroutine
	go func() {
		subject.Next(123)
		subject.Next(456)
		subject.Complete()
		subject.Next(888)
	}()

	// Subscribe to subject on the default Trampoline scheduler.
	subject.SubscribeNext(func(next int) {
		fmt.Println("first", next)
	})

	fmt.Println("--")

	// Subscribe to subject again on the default Trampoline scheduler.
	subject.SubscribeNext(func(next int) {
		fmt.Println("second", next)
	})

	// Output:
	// first 123
	// first 456
	// --
	// second 123
	// second 456
}

// ReplaySubject example with multiple subscribers. Subscribe normally uses a
// synchronous Trampoline scheduler. To be able to subscribe multiple times
// without blocking, we have to change the scheduler to an asynchronous one
// using the SubscribeOn operator.
func Example_replaySubjectMultiple() {
	subject := NewReplaySubjectString(1000, time.Hour)

	// Subscribe to subject on a goroutine
	source := subject.SubscribeOn(GoroutineScheduler())

	var results []*[]string
	var subscriptions []Subscription

	// Create a lot of subscribers running on separate goroutines.
	for i := 0; i < 16; i++ {
		result := make([]string, 0, 1000)
		results = append(results, &result)
		subscription := source.SubscribeNext(func(next string) {
			result = append(result, next)
		})
		subscriptions = append(subscriptions, subscription)
	}

	// Feed the subject from the main goroutine
	expect := make([]string, 0, 1000)
	for i := 0; i < 1000; i++ {
		s := fmt.Sprintf("Next %d", i)
		expect = append(expect, s)
		subject.Next(s)
	}
	subject.Complete()

	// Wait for all subscriptions to complete
	for _, subscription := range subscriptions {
		subscription.Wait()
	}

	// Check the results for missing values
	for index, result := range results {
		identical := false
		result := *result
		if len(result) == len(expect) {
			identical = true
			for index, value := range result {
				if value != expect[index] {
					identical = false
				}
			}
		}
		if identical {
			fmt.Printf("results %d as expected\n", index)
		}
	}

	// Output:
	// results 0 as expected
	// results 1 as expected
	// results 2 as expected
	// results 3 as expected
	// results 4 as expected
	// results 5 as expected
	// results 6 as expected
	// results 7 as expected
	// results 8 as expected
	// results 9 as expected
	// results 10 as expected
	// results 11 as expected
	// results 12 as expected
	// results 13 as expected
	// results 14 as expected
	// results 15 as expected
}
