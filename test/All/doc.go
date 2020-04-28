/*
	All	http://reactivex.io/documentation/operators/all.html

All determines whether all items emitted by an Observable meet some
criteria.

Pass a predicate function to the All operator that accepts an item emitted
by the source Observable and returns a boolean value based on an
evaluation of that item. All returns an ObservableBool that emits a single
boolean value: true if and only if the source Observable terminates
normally and every item emitted by the source Observable evaluated as
true according to the predicate; false if any item emitted by the source
Observable evaluates as false according to the predicate.
*/
package All

import _ "github.com/reactivego/rx"

// Don't generate documentation into rx.go
//off jig:no-doc
