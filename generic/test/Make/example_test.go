
package Make

import "time"

//jig:no-doc

func Example() {
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

func Example_interval() {
	const _10ms = 10 * time.Millisecond

	interval := func(period time.Duration) ObservableInt {
		i := 0
		return MakeTimedInt(period,
			func(Next func(int), Error func(error), Complete func()) time.Duration {
				Next(i)
				i++
				return period
			})
	}

	interval(_10ms).Take(4).Println()

	// Output:
	// 0
	// 1
	// 2
	// 3
}
