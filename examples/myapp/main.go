package main

import (
	"fmt"
	obs "rxgo/observable"
)

func main() {
	// import "..../observable"
	observable.CreateInt(func(subscriber observable.IntSubscriber) {

	})

	// import obs "..../observable"
	obs.CreateInt(func(subscriber obs.IntSubscriber) {

	})

	// import . "..../observable"
	CreateInt(func(subscriber IntSubscriber) {

	})

	a := obs.FromInts(1, 2, 3, 4).MapString(func(i int) string { return fmt.Sprintf("%d!", i) }).ToArray()
}
