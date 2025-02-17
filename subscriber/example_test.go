package subscriber

import  "fmt"

func Example_basic() {
	var s subscriber
	if s.Subscribed() == s.Canceled() {
		fmt.Println("subscriber can't be both subscribed and canceled")
	}
	if s.Canceled() {
		fmt.Println("subscriber should be subscribed by default")
	}
	s.Unsubscribe()
	if s.Subscribed() {
		fmt.Println("subscriber should be canceled after Unsubscribe call")
	}

	fmt.Println("OK")
	// Output: OK
}


func Example_subscriberLoop() {
	parent := &subscriber{}

	child1 := parent.Add()
	child2 := parent.Add()
	child2.OnUnsubscribe(parent.Unsubscribe)
	child3 := parent.Add()

	if parent.Canceled() || child1.Canceled() || child2.Canceled() || child3.Canceled() {
		fmt.Println("none of the subscribers should be cancelled here")
	}

	child2.Unsubscribe()

	if parent.Subscribed() || child1.Subscribed() || child2.Subscribed() || child3.Subscribed() {
		fmt.Println("all of the subscribers should be cancelled here")
	}

	fmt.Println("OK")
	// Output: OK
}