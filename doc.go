// Package rx offers Reactive Extensions (ReactiveX) for Go, an API for
// asynchronous programming with observable streams.
//
// The coherent set of functions and types offered by this package make it
// easier to develop event-based and asynchronous programs using observable
// streams (Observables).
//
// Heterogeneous vs typesafe
// 
// A sequence or stream can be either:
//	typesafe, guaranteed to contain items of a single type. e.g. string.
//	heterogeneous, contain items of any type.
//
// The implementation in this package is rx/generic specialized on interface{}.
// As a value of interface{} accepts any type, this makes this package heterogeneous.
//
// By using Generics for Go (https://github.com/reactivego/jig) it is also
// possible to have the jig command generate typesafe implementations.
//
// Alias any
//
// In this package we have defined a non-exported alias 'any' for 'interface{}'.
// So in the examples, instead of 'interface{}' you will find 'any' being used.
// If you choose to use this package in your own code, you can decide for yourself
// how to call the interface{} type, this is not enforced by this package.
//
// Creating Observables
//
// Functions to create new observables from lists, slices and channels.
//	From          http://reactivex.io/documentation/operators/from.html
//	FromChan
//	FromSlice
//
// Combining Observables
//
// Functions and methods on observable that allow combining multiple observables.
// 	Concat        http://reactivex.io/documentation/operators/concat.html
//
// Transforming Observables
//
// Functions and methods to transform observables to same or other types.
// 	Filter        http://reactivex.io/documentation/operators/filter.html
// 	MergeMap      http://reactivex.io/documentation/operators/flatmap.html
// 	Map           http://reactivex.io/documentation/operators/map.html
// 	Scan          http://reactivex.io/documentation/operators/scan.html
//
// Side Effects triggered by Observables
//
// Methods on observables to trigger side effects.
//	Do            http://reactivex.io/documentation/operators/do.html
//	DoOnError
//	DoOnComplete
//	Finally
//
// For an overview of all implemented operators, see
// https://github.com/reactivego/rx
package rx

import _ "github.com/reactivego/rx/generic"

type any = interface{}
