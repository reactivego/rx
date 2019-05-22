package Only

import "fmt"

// Size is recognized by jig directly because refer type and actual type are both "Size"
type Size struct{ width, height float64 }

// point will be used in a []point
type point struct{ x, y float64 }

// Point is how we refer to the type []point. Jig must be told since it can't know that.
//jig:type Point []point

func Example_only() {
	// Create an Observable of interface{} and stuff it with values of different types.
	source := Create(func(observer Observer) {
		observer.Next("Hello")                      // string
		observer.Next(Size{1024, 768})              // Size
		observer.Next([]point{{50, 100}, {75, 25}}) // []point
		observer.Complete()
	})
	// OnlyString() will filter on string values
	source.OnlyString().SubscribeNext(func(next string) { fmt.Printf("String: %s\n", next) })
	// OnlySize() will filter on Size values
	source.OnlySize().SubscribeNext(func(next Size) { fmt.Printf("Size: %+v\n", next) })
	// OnlyPoint() will filter on []point values
	source.OnlyPoint().SubscribeNext(func(next []point) { fmt.Printf("Point: %+v\n", next) })

	// Output:
	// String: Hello
	// Size: {width:1024 height:768}
	// Point: [{x:50 y:100} {x:75 y:25}]
}
