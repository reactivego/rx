package Make

import _ "github.com/reactivego/rx"

func Example_makeInt() {
	example := func() ObservableInt {
		done := false
		return MakeInt(func(Next func(int), Error func(error), Complete func()) {
			if !done {
				Next(1)
				done = true
			} else {
				Complete()
			}
		})
	}

	example().Println()

	// Output:
	// 1
}
