package main

// For this example to work you need to first run rxgo to generate an appropriate rx.go file alongside this main.go file.
// e.g. rxgo main int --output=rx.go --map-types=MouseMove
// The generated file will have Map operator from observable int to observable MouseMove

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
