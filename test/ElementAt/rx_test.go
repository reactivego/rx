package ElementAt

import _ "github.com/reactivego/rx"

func Example_elementAt() {
	FromInt(1, 2, 3, 4).ElementAt(2).Println()
	// Output:
	// 3
}
