package rx

import "fmt"

func ExampleObservable_Do() {
	From(1,2,3).Do(func(v any) {
		fmt.Println(v.(int))
	}).Wait()

	// Output:
	// 1
	// 2
	// 3
}
