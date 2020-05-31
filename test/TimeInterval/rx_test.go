package TimeInterval

import (
	"time"

	_ "github.com/reactivego/rx"
)

func Example_timeInterval() {
	const ms = time.Millisecond

	Timer(10*ms, 100*ms).Take(5).TimeInterval().Println()
	// Output:
	// {0 10ms}
	// {1 100ms}
	// {2 100ms}
	// {3 100ms}
	// {4 100ms}
}
