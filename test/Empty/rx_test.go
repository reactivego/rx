package Empty

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_empty() {
	err := Empty().Println()
	fmt.Println(err)
	// Output: <nil>
}
