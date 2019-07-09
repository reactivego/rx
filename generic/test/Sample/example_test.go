package Sample

import (
	"fmt"
	"time"
)

func Example_sample() {
	const _90ms = 90 * time.Millisecond
	const _200ms = 200 * time.Millisecond
	obsrvbl := Interval(_90ms).Sample(_200ms).Take(3)
	if err := obsrvbl.Println(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// 1
	// 3
	// 5
}
