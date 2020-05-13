// Package rx/generic provides Reactive Extensions for Go.
// A generics library for asynchronous programming with observable streams.
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
