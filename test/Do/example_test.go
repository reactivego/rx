package Do

import "fmt"

func Example_do() {
	slice := []int(nil)

	FromInts(1, 2, 3, 4, 5).Do(func(next int) {
		slice = append(slice, next)
	}).Wait()

	fmt.Println(slice)

	// Output: [1 2 3 4 5]
}
