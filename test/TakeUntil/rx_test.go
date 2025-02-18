package TakeUntil

import (
	"time"

	_ "github.com/reactivego/rx/generic"
)

func Example_takeUntil() {
	const ms = time.Millisecond

	// produce a single value after 550ms
	interrupted := Timer(550 * ms)

	// produce a number starting at 0ms every 100ms
	source := Timer(0*ms, 100*ms)

	// print numbers from source until interrupted
	source.TakeUntil(interrupted).Println()

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
}
