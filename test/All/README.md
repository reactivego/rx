# All

[![](../../svg/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx/test/All?tab=doc)
[![](../../svg/godoc.svg)](https://godoc.org/github.com/reactivego/rx/test/All)
[![](../../svg/rx.svg)](http://reactivex.io/documentation/operators/all.html)

**All** determines whether all items emitted by an Observable meet some
criteria.

Pass a predicate function to the **All** operator that accepts an item emitted
by the source Observable and returns a boolean value based on an
evaluation of that item. **All** returns an ObservableBool that emits a single
boolean value: true if and only if the source Observable terminates
normally and every item emitted by the source Observable evaluated as
true according to the predicate; false if any item emitted by the source
Observable evaluates as false according to the predicate.

<!--
marble all
{
	source a:                 +-1-2-6-2-1-|
	operator All(Less<Than5): +-----------(false)|
}
-->
![All](../../svg/All.svg)
## Example
```go
import _ "github.com/reactivego/rx"
```

Code:
```go
LessThan5 := func(i int) bool {
	return i < 5
}
FromInt(1, 2, 6, 2, 1).All(LessThan5).Println("All values less than 5?")
```

Output:
```
All values less than 5? false
```
