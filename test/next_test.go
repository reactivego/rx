package test

import (
	"errors"
)

// Next structs are used in ToChanNext to report errors in-band along with the
// data stream. This example shows how to create a channel of NextInt with a
// buffer of 2 and how to send the next value or an error.
func Example_next() {
	ch := make(chan NextInt, 2)
	defer close(ch)

	// Send Next Value
	ch <- NextInt{Next: 123}

	// Send Error
	ch <- NextInt{Err: errors.New("error")}

	// Output:
}
