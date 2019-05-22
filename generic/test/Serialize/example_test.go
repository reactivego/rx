package Serialize

import (
	"fmt"
	"sync/atomic"
	"time"
)

// In this example we deliberately create a badly behaved observable that
// violates the observable contract. We verify that we can detect that, and then
// we add Serialize and show that the contract is now correctly enforced.
func Example_serialize() {
	// Create badly behaved observable.
	source := CreateInt(func(observer IntObserver) {
		hammer := func(value int) {
			for !observer.Closed() {
				observer.Next(value)
			}
		}
		go hammer(1)
		go hammer(2)
		go hammer(3)
	})

	// Detecting violation of the observable contract by counting how many
	// times multiple goroutines got into the observer at the same time.
	var concurrent, violations int32
	violationChecker := func(next int) {
		if atomic.AddInt32(&concurrent, 1) > 1 {
			atomic.AddInt32(&violations, 1)
		}
		atomic.AddInt32(&concurrent, -1)
	}

	// Subscribe the detector on the source directly and let it 'hammer'
	// the observer for a millisecond.
	subscription := source.SubscribeNext(violationChecker)
	time.Sleep(time.Millisecond)
	subscription.Unsubscribe()
	fmt.Printf("Observable contract violated = %v\n", violations > 0)

	// Reset the violation checker
	concurrent, violations = 0, 0

	// Serialize the source.
	source = source.Serialize()

	// Subscribe the detector on the source directly and let it 'hammer'
	// the observer for a millisecond.
	subscription = source.SubscribeNext(violationChecker)
	time.Sleep(time.Millisecond)
	subscription.Unsubscribe()
	fmt.Printf("Observable contract violated = %v\n", violations > 0)

	// Output:
	// Observable contract violated = true
	// Observable contract violated = false
}
