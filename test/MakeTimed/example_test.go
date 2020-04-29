package MakeTimed

import (
	"time"

	_ "github.com/reactivego/rx"
)

func Example_makeTimedInt() {
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
