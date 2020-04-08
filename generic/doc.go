/*
Package rx provides generic Reactive Extensions (ReactiveX) for Go.
It's a generics library for composing asynchronous and event-based programs
using observable sequences. The library consists of more than a 100 generic 
templates to enable type-safe programming with observable streams.
To use it, you will need Generics for Go (https://github.com/reactivego/jig).

Using the library is simple. Start by creating a file "main.go":

	package main

	import _ "github.com/reactivego/rx/generic"

	func main() {
		FromInt(1,2,3).Println()
	}
	
The import is purely done for jig. Jig searches for generic templates inside
the imported package. FromInt matches generic function template From<Foo>.
Running jig will generate the code:

	go get github.com/reactivego/rx
	go get github.com/reactivego/jig
	go run github.com/reactivego/jig -rv

The code has been generated into file "rx.go".
You can now build and run the generated code:

	go run *.go

Take a look at the Quick Start guide to see a more elaborate version of the
example above. https://github.com/ReactiveGo/jig/blob/master/QUICKSTART.md.

For an overview of all implemented operators, see
https://github.com/reactivego/rx/tree/master/generic#operators
*/
package rx

// foo is the first metasyntactic type. Use the jig:type pragma to tell jig that
// Foo is the reference type name for actual type foo. Needed because we're
// generating code into rx.go for foo.

//jig:type Foo foo

type foo int

// bar is the second metasyntactic type. Use the jig:type pragma to tell jig that
// Bar is the reference type name for actual type bar. Needed because we're
// generating code into rx.go for bar.

//jig:type Bar bar

type bar int
