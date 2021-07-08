package BufferTime

import (
	"time"

	_ "github.com/reactivego/rx"
)

func Example_bufferTime() {
	const ms = time.Millisecond
	source := TimerInt(0*ms, 100*ms).Take(4).ConcatMap(func(i int) Observable {
		switch i {
		case 0:
			return From("a", "b")
		case 1:
			return From("c", "d", "e")
		case 3:
			return From("f", "g")
		}
		return Empty()
	})
	source.BufferTime(100 * ms).Println()

	// Output:
	// [a b]
	// [c d e]
	// []
	// [f g]
}
