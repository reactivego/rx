package test

import (
	"fmt"
	"sync"
)

// ChanNext is a struct wrapped around a channel of Next. It has methods to
// Send data and to Close and Cancel the channel. Cancel is meant to be called
// from the receiving side.
func Example_chanNext() {
	ch := NewChanNextInt(2)

	for i := 123; i < 125; i++ {
		ch.Send(i)
		fmt.Printf("sent %d\n", i)
	}

	fmt.Println((<-ch.Channel).Next)
	fmt.Println((<-ch.Channel).Next)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		ch.Send(111)
		ch.Send(222)
		ch.Send(333)
		ch.Send(444)
		ch.Send(555)
		ch.Send(666)
		ch.Send(777)
		ch.Send(888)
		ch.Send(999)
		fmt.Println("finished")
		wg.Done()
	}()

	// Wait for the sender goroutine to start sending.
	fmt.Println((<-ch.Channel).Next)
	// Channel sender will now be blocked on ch.Send(444)

	// Tell the channel to cancel.
	ch.Cancel()

	// Wait for the canceled sender goroutine to finish.
	wg.Wait()

	// Will print 222 and 333 that were in channel before Cancel.
	for entry := range ch.Channel {
		fmt.Println(entry.Next)
	}

	// Output:
	// sent 123
	// sent 124
	// 123
	// 124
	// 111
	// finished
	// 222
	// 333
}
