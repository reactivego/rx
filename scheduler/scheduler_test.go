package scheduler

import "fmt"

/*
// The serial queue abstraction may be a good fit for supplying concurrency into the subscribe and observe actions.
func main() {
	q := StartSerialQueue()


	// this will deadlock
	q.Do(....)

	// always use
	q.Go(...
		q.Do(....)
		...
	)

	// Range will schedule its mainloop function via q.Go while ObserveOn will schedule the
	// delivery of the data on the main thread.
	// Note how the statement below creates a reusable component with concurrency parameterized into it.
	r := Range(1, 100).SubscribeOn(q.Go).ObserveOn(q.Do)

	// First computation, running on a separate goroutine via q.Go
	s := r.Subscribe(func(next int, err error, completed bool) {
		// This function called via q.Do will run on the main thread

		// You can assign data to the global variables here, really!
	})

	// Second computation, running on a separate goroutine via q.Go
	s := r.Subscribe(func(next int, err error, completed bool) {
		// This function called via q.Do will run on the main thread

		// You can assign data to the global variables here, really!
	})

	// The two computations above will run concurrently and will be
	// serialized. Note that the second computation may finish before
	// the first computation. The guarantee is that the


	// because this will deadlock on ToArray we need to wrap this in a Go call.
	q.Go(func() {
		a := r.ToArray()
		q.Do(func() {
			// you could assign a to a global variable here, as this is running on the main thread
		})
	})

	// Wait for all the q.Go() functions to complete.
	q.Wait()


	// when we get here all computations have finished.
}
*/

func ExampleSimpleGo1() {

	// Running just go routine will never print because program has already terminated
	go func() {
		fmt.Println("Hello, World!")
	}()

	//Output:
}

func ExampleSimpleGo2() {
	wait := make(chan struct{})
	go func() {
		defer close(wait)
		fmt.Println("Hello, World!")
	}()
	<-wait
	//Output:
	// Hello, World!
}

func ExampleSimpleSerialQueue() {
	q := StartSerialQueue()
	q.Go(func() {
		fmt.Println("Hello, World!")
	})
	q.Wait()
	//Output:
	// Hello, World!
}

func ExampleNestedSerialQueue() {
	q := StartSerialQueue()
	q.Go(func() {
		fmt.Println("Hello, World!")

		// start a computation and
		r := StartSerialQueue()
		r.Go(func() {
			// do lots of computation
			r.Do(func() {
				// will be scheduled on the goroutine in which r.Wait() is called.
				fmt.Println("Hello, Again!")
			})
		})

		// start more concurrent work using r.Go() and deliver results to the parent goroutine via r.Do()

		// Wait for all r.Go() code to terminate.
		r.Wait()

		// use the results produced by the computations that ran via r.Go()
		fmt.Println("Yes, Hello!")
	})
	q.Wait()
	//Output:
	// Hello, World!
	// Hello, Again!
	// Yes, Hello!
}

func Example() {
	q := StartSerialQueue()

	// Code running here is on the main goroutine.

	q.Go(func() {
		// Code executed here runs on a separate 'background work' goroutine.
		fmt.Println("Background Work A")
		// func called via Do runs on the main goroutine!
		q.Do(func() { fmt.Println("Deliver Work A") })

		fmt.Println("Background Work B")
		q.Do(func() { fmt.Println("Deliver Work B") })

		fmt.Println("Background Work C")
		q.Do(func() { fmt.Println("Deliver Work C") })

		fmt.Println("Background Work D")
		q.Do(func() { fmt.Println("Deliver Work D") })

		fmt.Println("Background Work E")
		q.Do(func() { fmt.Println("Deliver Work E") })

		fmt.Println("Background Work F")
		q.Do(func() { fmt.Println("Deliver Work F") })

		fmt.Println("Background Work G")
		q.Do(func() { fmt.Println("Deliver Work G") })

		fmt.Println("Background Work H")
		q.Do(func() { fmt.Println("Deliver Work H") })
	})

	q.Wait()

	//Output:
	// Background Work A
	// Deliver Work A
	// Background Work B
	// Deliver Work B
	// Background Work C
	// Deliver Work C
	// Background Work D
	// Deliver Work D
	// Background Work E
	// Deliver Work E
	// Background Work F
	// Deliver Work F
	// Background Work G
	// Deliver Work G
	// Background Work H
	// Deliver Work H
}
