package main

import _ "github.com/reactivego/rx"

//jig:file rx.go

func main() {
	FromStrings("You!", "Gophers!", "World!").
		MapString(func(x string) string {
			return "Hello, " + x
		}).
		SubscribeNext(func(next string) {
			println(next)
		})

	// Output:
	// Hello, You!
	// Hello, Gophers!
	// Hello, World!
}
