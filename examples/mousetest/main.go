package main

import (
	"rxgo/examples/mousetest/mouse"
	"rxgo/observable"
)

func main() {
	seq := observable.FromIntArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	moves := seq.MapMove(func(v int) *mouse.Move {
		return &mouse.Move{X: float64(v), Y: float64(v)}
	})

	done := make(chan struct{})
	moves.Subscribe(func(next *mouse.Move, err error, complete bool) {
		if err == nil && !complete {
			println(next.X, next.Y)
		} else {
			close(done)
		}
	})
	<-done
}
