package All

import "fmt"

func Example_all() {
	source := CreateInt(func(observer IntObserver) {
		observer.Next(1)
		observer.Next(2)
		observer.Next(6)
		observer.Next(2)
		observer.Next(1)
		observer.Complete()
	})

	// Print the sequence of numbers
	source.Println()
	fmt.Println("Source observable completed")

	// Setup All to produce true only when all source values are less than 5
	predicate := func(i int) bool {
		return i < 5
	}

	// Observer prints a message describing the next value it observes.
	observer := func(next bool, err error, done bool) {
		if !done {
			fmt.Println("All values less than 5?", next)
		}
	}

	source.All(predicate).Subscribe(observer).Wait()
	// Output:
	// 1
	// 2
	// 6
	// 2
	// 1
	// Source observable completed
	// All values less than 5? false
}
