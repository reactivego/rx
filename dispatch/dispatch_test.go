package dispatch

import (
	"fmt"
	"testing"
)

func TestCounterAccess(t *testing.T) {
	operations := NewSerialQueue()
	operations.Start()

	// Note that the loop below will launch 1000 concurrent goroutines
	// via the DispatchAsync call. However the actual code accessing the
	// counter will all be running in the single goroutine associated with
	// the operations queue

	var counter int
	for i := 0; i < 1000; i++ {
		// force value copy of i in c
		c := i
		operations.DispatchAsync(func() {
			// this function is performed by a single goroutine associated with the operations queue
			counter += c
		})
	}
	operations.Stop()
	if counter != 499500 {
		t.Errorf("counter should be 499500 is %d", counter)
	}
}

func ExampleMainSerialQueue() {
	fmt.Println("hello from main.main: BEGIN")
	Main.Go(func(q *SerialQueue) {
		fmt.Println("hello from goroutine: BEGIN")
		// tasks dispatched here are guaranteed to be performed by the queue
		q.DispatchSync(func() { fmt.Println("hello from synchronous task...(1)") })
		q.DispatchSync(func() { fmt.Println("hello from synchronous task...(2)") })
		q.DispatchSync(func() { fmt.Println("hello from synchronous task...(3)") })
		// statement below is guaranteed to print.
		fmt.Println("hello from goroutine: END")
	})

	Main.Main()
	fmt.Println("hello from main.main: END")

	//Output:
	// hello from main.main: BEGIN
	// hello from goroutine: BEGIN
	// hello from synchronous task...(1)
	// hello from synchronous task...(2)
	// hello from synchronous task...(3)
	// hello from goroutine: END
	// hello from main.main: END
}
