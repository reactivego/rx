package CreateRecursive

import _ "github.com/reactivego/rx"

func Example_createRecursiveInt() {
	example := func() ObservableInt {
		done := false
		return CreateRecursiveInt(func(N NextInt, E Error, C Complete) {
			if !done {
				N(1)
				done = true
			} else {
				C()
			}
		})
	}

	example().Println()

	// Output:
	// 1
}
