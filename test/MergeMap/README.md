# MergeMap

[![](../../svg/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx/test/MergeMap?tab=doc)
[![](../../svg/godoc.svg)](http://godoc.org/github.com/reactivego/rx/test/MergeMap)
[![](../../svg/rx.svg)](http://reactivex.io/documentation/operators/flatmap.html)

**MergeMap** transforms the items emitted by an Observable by applying a
function to each item an returning an Observable. The stream of Observable
items is then merged into a single stream of items using the MergeAll
operator.

This operator was previously named FlatMap. The name FlatMap is deprecated as
MergeMap more accurately describes what the operator does with the observables
returned from the Map project function.
