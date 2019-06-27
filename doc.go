/*
Package rx provides Reactive Extensions (ReactiveX) for Go, an API for
asynchronous programming with observable streams (Observables).

Heterogeneous

The implementation in this package is created by expanding templates from
"github.com/reactivego/rx/generic" on the "interface{}" type. A variable of 
type "interface{}"" accepts a value of any type. This makes this package
heterogeneous, since the definition of a heterogeneous variable is that it
can contain a value of any type. 

	From(1,"hi",2.3).Println()

Heterogenous is great when you want to mix values of different types in a
stream, but by gaining that, you loose type safety. However, by using
Generics for Go (https://github.com/reactivego/jig) it is also possible to
have the jig command generate typesafe implementation of the templates.

	FromInt(1,2,3).Println()

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
https://github.com/reactivego/rx
*/
package rx

import _ "github.com/reactivego/rx/generic"

// In this package we have defined a non-exported alias 'any' for 'interface{}'.
// So in the examples, instead of 'interface{}' you will find 'any' being used.
// If you choose to use this package in your own code, you can decide for yourself
// how to call the interface{} type, this is not enforced by this package.
type any = interface{}
