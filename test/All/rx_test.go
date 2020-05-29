package All

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_all() {
	source := CreateInt(func(N NextInt, E Error, C Complete, X Canceled) {
		N(1)
		N(2)
		N(6)
		N(2)
		N(1)
		C()
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
