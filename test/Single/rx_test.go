package Single

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

// Single is used to make sure only a single value was produced by the
// observable.
func Example_single() {
	// Just output 1 int.
	FromInt(19).Single().Subscribe(func(next int, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}).Wait()

	// Now output 2 ints.
	FromInt(19, 20).Single().Subscribe(func(next int, err error, done bool) {
		switch {
		case !done:
			fmt.Println(next)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}).Wait()
	// Output:
	// 19
	// complete
	// expected one value, got multiple
}
