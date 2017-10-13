package channel

import "fmt"

// Example uses all the functionality that we want to export for this channel
// package.
func Example() {
	ch := NewChan(128, 1)
	if true {
		ch.FastSend("hello")
	} else {
		ch.Send("world!")
	}
	ch.Close(nil)
	if ch.Closed() {
		fmt.Println("closed")
	}

	ep, _ := ch.NewEndpoint(ReplayAll)
	print := func(value interface{}, err error, closed bool) bool {
		switch {
		case !closed:
			fmt.Println(value)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("closed")
		}
		return true
	}
	ep.Range(print, 0)
	ep.Cancel()

	// Output:
	// closed
	// hello
	// closed
}
