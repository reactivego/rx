// Library rx provides a (jig) generic implementation of Reactive eXtensions.
// It's a generics library for composing asynchronous and event-based programs
// using observable sequences. The library consists of more than a 100 templates
// to enable type-safe programming with observable streams. To use it, you will
// need the jig tool https://github.com/reactivego/jig.
//
// Using the library is very simple. Import the library with the blank
// identifier `_` as the package name:
//	import _ "github.com/reactivego/rx"
// The side effect of this import is that
// generics from the library can now be accessed by the jig tool. Then start
// using generics from the library and run jig to generate code. Take a look at
// the Quick Start guide to see how it all fits together
// https://github.com/reactivego/rx/doc/quickstart.md.
//
// For an overview of all implemented operators, see
// https://github.com/reactivego/rx
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
