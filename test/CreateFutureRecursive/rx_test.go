package CreateFutureRecursive

import (
	"time"

	_ "github.com/reactivego/rx/generic"
)

func Example_createFutureRecursiveInt() {
	const _10ms = 10 * time.Millisecond

	interval := func(period time.Duration) ObservableInt {
		i := 0
		return CreateFutureRecursiveInt(period, func(N NextInt, E Error, C Complete) time.Duration {
			N(i)
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
