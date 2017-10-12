// Library rx provides (jig) Reactive eXtensions for the Go language.
//
// This is a 'Just-In-time Generics' library for Go. It can be used for
// composing asynchronous and event-based programs using observable sequences.
// This package contains more than a 100 templates to enable type-safe
// programming with observable streams.
//
// Because the generics definitions in 'rx' are only recognized by
// Just-in-time Generics for Go, you will need to install the jig tool
// (https://github.com/reactivego/jig/). Don’t worry, it's easy to
// setup and use, take a look at the Quick Start guide.
// (https://github.com/reactivego/rx/doc/quickstart.md).
//
// For an overview of all implemented operators, see
// https://github.com/reactivego/rx
package rx

// Code generated to make this rx package buildable should be written to rx.go

//jig:file {{.package}}.go

// Support code is code that is not parameterized on a template variable. Code
// like that (when generated) is the same as the code already present in this
// package and would therefore lead to symbol collisions. Disable generation
// of support code using the jig:no-support-code-generation pragma.

//jig:no-support-code-generation

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
