package main

type MouseMove struct {
	X float64
	Y float64
}

func main() {
	seq := FromIntArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	moves := seq.MapMouseMove(func(v int) MouseMove {
		return MouseMove{X: float64(v), Y: float64(v)}
	})

	done := make(chan struct{})
	moves.Subscribe(func(next MouseMove, err error, complete bool) {
		if err == nil && !complete {
			println(next.X, next.Y)
		} else {
			close(done)
		}
	})
	<-done
}
