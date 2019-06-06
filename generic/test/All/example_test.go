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

	// Setup All operator to produce true only when all source values are less than 5
	all := source.All(func(i int) bool { return i < 5 })

	all.SubscribeNext(func(b bool) { fmt.Println("All values less than 5?", b) }).Wait()

	fmt.Println("All operator completed")

	// Output:
	// 1
	// 2
	// 6
	// 2
	// 1
	// Source observable completed
	// All values less than 5? false
	// All operator completed
}
