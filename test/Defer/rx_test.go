package Defer

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_defer() {
	count := 0
	source := DeferInt(func() ObservableInt {
		// called for every subscribe of an observer
		count++
		return FromInt(count)
	})
	mapped := source.MapString(func(next int) string {
		return fmt.Sprintf("observer %d", next)
	})

	// subscribe first observer
	mapped.Println()

	// subscribe second observer
	mapped.Println()
	// Output:
	// observer 1
	// observer 2
}
