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

	fobservable := observable.Distinct().MapFloat64(func(v int) float64 { return float64(v) * 1.62 })

	term := make(chan struct{})
	subscription := fobservable.Subscribe(func(v float64, err error, complete bool) {
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
