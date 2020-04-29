package Sample

import (
	"fmt"
	"time"

	_ "github.com/reactivego/rx"
)

func Example_sample() {
	const _90ms = 90 * time.Millisecond
	const _200ms = 200 * time.Millisecond
	err := Interval(_90ms).Sample(_200ms).Take(3).Println()
	fmt.Println(err)
	// Output:
	// 1
	// 3
	// 5
	// <nil>
}
