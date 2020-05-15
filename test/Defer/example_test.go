package Defer

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_defer() {
	count := 0
	source := DeferInt(func() ObservableInt {
		count++
		return FromInt(count)
	})
	mapped := source.MapString(func(next int) string {
		return fmt.Sprintf("observable %d", next)
	})

	mapped.Println()
	mapped.Println()
	// Output:
	// observable 1
	// observable 2
}
