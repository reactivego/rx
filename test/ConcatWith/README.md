# ConcatWith

[![](../../svg/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx/test/ConcatWith?tab=doc)
[![](../../svg/godoc.svg)](http://godoc.org/github.com/reactivego/rx/test/ConcatWith)
[![](../../svg/rx.svg)](http://reactivex.io/documentation/operators/concat.html)

**ConcatWith** emits the emissions from two or more observables without interleaving them.

<!--
marble ConcatWith
{
	source a:            +-1----1------1--|
	source b:            +-2-2-|
	operator ConcatWith: +-1----1------1---2-2-|
}
-->
![ConcatWith](../../svg/ConcatWith.svg)

## Example
```go
import _ "github.com/reactivego/rx"
```
Code:
```go
oa := FromInt(0, 1, 2)
ob := FromInt(3)
oc := FromInt(4, 5)
oa.ConcatWith(ob, oc).Println()
```
Output:
```
0
1
2
3
4
5
```
