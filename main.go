package main

import "rxgo/observable"

////////////////////////////////////////////////////////
// main
////////////////////////////////////////////////////////

func main() {
	println("hello")

	// observable := Range(1, 10)
	observable := observable.CreateInt(func(s observable.IntSubscriber) {
		println("subscription received...")
		for i := 0; i < 15; i++ {
			for j := 0; j < 3; j++ {
				if s.Disposed() {
					return
				}
				s.Next(i)
			}
		}
		s.Complete()
		s.Dispose()
	})

	observable = observable.Distinct()

	term := make(chan struct{})
	subscription := observable.SubscribeFunc(func(v int, err error, complete bool) {
		if err != nil || complete {
			close(term)
			return
		}
		println(v)
	})
	<-term
	if subscription.Disposed() {
		println("subscription disposed....")
	} else {
		println("subscription still alive....")
	}
	println("goodbye")
}
