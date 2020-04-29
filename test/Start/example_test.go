package Start

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_start() {
	StartInt(func() (int, error) { return 42, nil }).Println()
	// Output: 42
}

func Example_startWithError() {
	err := StartInt(func() (int, error) { return 0, RxError("start-with-error") }).Println()
	fmt.Println(err)
	// Output:
	// start-with-error
}
