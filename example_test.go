package subscriber

import  "fmt"

func Example_basic() {
	var s subscriber
	if s.Closed() {
		fmt.Println("subscriber should not be closed by default")
	}
	s.Unsubscribe()
	if !s.Closed() {
		fmt.Println("subscriber should be closed after Unsubscribe call")
	}

	fmt.Println("OK")
	// Output: OK
}


func Example_subscriberLoop() {
	parent := &subscriber{}

	child1 := parent.Add()
	child2 := parent.Add(parent.Unsubscribe)
	child3 := parent.Add()

	if parent.Canceled() || child1.Canceled() || child2.Canceled() || child3.Canceled() {
		fmt.Println("none of the subscribers should be cancelled here")
	}

	child2.Unsubscribe()

	if !(parent.Canceled() && child1.Canceled() && child2.Canceled() && child3.Canceled()) {
		fmt.Println("all of the subscribers should be cancelled here")
	}

	fmt.Println("OK")
	// Output: OK
}