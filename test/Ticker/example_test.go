package Ticker

import (
	"time"

	_ "github.com/reactivego/rx"
)

func Example_ticker() {
	const ms = time.Millisecond

	project := func(t time.Time) string {
		return t.Format("2006")
	}
	Ticker(10*ms, 100*ms).MapString(project).Take(5).Println()
	// Output:
	// 2020
	// 2020
	// 2020
	// 2020
	// 2020
}
