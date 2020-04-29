package DoOnError

import (
	"fmt"

	_ "github.com/reactivego/rx"
)

func Example_doOnError() {
	var oerr error
	err := Throw(RxError("bazinga!")).DoOnError(func(err error) { oerr = err }).Wait()

	fmt.Println(oerr, err)
	// Output: bazinga! bazinga!
}
