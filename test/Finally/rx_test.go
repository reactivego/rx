package Finally

import (
	"fmt"

	_ "github.com/reactivego/rx/generic"
)

func Example_basic() {
	err := ThrowInt(RxError("error")).Finally(func() { fmt.Println("finally") }).Wait()
	fmt.Println(err)

	err = EmptyInt().Finally(func() { fmt.Println("finally") }).Wait()
	fmt.Println(err)

	// Output:
	// finally
	// error
	// finally
	// <nil>
}
