/*
Package rx provides Reactive Extensions (ReactiveX) for Go, an API for
asynchronous programming with observable streams (Observables).

What rx is

Observables are the main focus of rx and they are (sort of) sets of events.
They assume zero to many values over time.
They push values.
They can take any amount of time to complete (or may never).
They are cancellable.
They are lazy; they don't do anything until you subscribe.

This implementation of rx uses interface{} for value types, so you can
mix different types of values in function and method calls. To create an
observable you might write the following:

	package main

	import "github.com/reactivego/rx"

	func main() {
		rx.From(1,"hi",2.3).Println()
	}

The code above creates an observable from numbers and strings and then prints
them.

Creating Observables

Functions to create new observables from lists, slices and channels.
	From          http://reactivex.io/documentation/operators/from.html
	FromChan
	FromSlice

Combining Observables

Functions and methods on observable that allow combining multiple observables.
	Concat        http://reactivex.io/documentation/operators/concat.html

Transforming Observables

Functions and methods to transform observables to same or other types.
	Filter        http://reactivex.io/documentation/operators/filter.html
	MergeMap      http://reactivex.io/documentation/operators/flatmap.html
	Map           http://reactivex.io/documentation/operators/map.html
	Scan          http://reactivex.io/documentation/operators/scan.html

Side Effects triggered by Observables

Methods on observables to trigger side effects.
	Do            http://reactivex.io/documentation/operators/do.html
	DoOnError
	DoOnComplete
	Finally

For an overview of all implemented operators, see
https://github.com/reactivego/rx/tree/master/generic#operators

Regenerating this Package

This package is generated from a generic implementation in the
subdirectory "generic" by a generator called "jig". You don't need "jig"
in order to use this package though. If you are however interested in
regenerating this package, then read on. Jig works by replacing the
place-holder types of templated types with interface{}. To regenerate
this rx implementation, run jig inside this package directory as follows:

	go get -d github.com/reactivego/generics/cmd/jig
	go run github.com/reactivego/generics/cmd/jig -v
*/
package rx

import _ "github.com/reactivego/rx/generic"

// In this package we have defined a non-exported alias 'any' for 'interface{}'.
// So in the examples, instead of 'interface{}' you will find 'any' being used.
// If you choose to use this package in your own code, you can decide for yourself
// how to call the interface{} type, this is not enforced by this package.
type any = interface{}
