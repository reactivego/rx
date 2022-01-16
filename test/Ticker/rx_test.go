package Ticker

import (
	"time"

	_ "github.com/reactivego/rx"
)

func Example_ticker() {
	const ms = time.Millisecond

	year := time.Now().Format("2006")

	project := func(t time.Time) string {
		if year == t.Format("2006") {
			return "OK"
		} else {
			return "FAILURE"
		}
	}
	Ticker(10*ms, 100*ms).MapString(project).Take(5).Println()
	// Output:
	// OK
	// OK
	// OK
	// OK
	// OK
}
