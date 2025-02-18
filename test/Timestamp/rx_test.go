package Timestamp

import (
	"time"

	_ "github.com/reactivego/rx/generic"
)

func Example_timestamp() {
	const ms = time.Millisecond

	check := func(t TimestampTime) bool {
		return t.Value.Round(ms) == t.Timestamp.Round(ms)
	}

	Ticker(10*ms, 100*ms).Timestamp().Take(5).All(check).Println("Timestamps correct =")
	// Output:
	// Timestamps correct = true
}
