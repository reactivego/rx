package Only

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

// Point is how we refer to the type []point. Jig must be told since it can't know that.
//jig:type Point []point

func Example_only() {
	// Create an Observable of interface{} and stuff it with values of different types.
	source := Create(func(N Next, E Error, C Complete, X Canceled) {
		N("Hello")                      // string
		N(Size{1024, 768})              // Size
		N([]point{{50, 100}, {75, 25}}) // []point
		C()
	})

	subscription := source.OnlyString().Subscribe(func(next string, err error, done bool) {
		if !done {
			fmt.Printf("String: %s\n", next)
		}
	})
	subscription.Wait()

	subscription = source.OnlySize().Subscribe(func(next Size, err error, done bool) {
		if !done {
			fmt.Printf("Size: %+v\n", next)
		}
	})
	subscription.Wait()

	subscription = source.OnlyPoint().Subscribe(func(next []point, err error, done bool) {
		if !done {
			fmt.Printf("Point: %+v\n", next)
		}
	})
	subscription.Wait()
	// Output:
	// String: Hello
	// Size: {width:1024 height:768}
	// Point: [{x:50 y:100} {x:75 y:25}]
}
