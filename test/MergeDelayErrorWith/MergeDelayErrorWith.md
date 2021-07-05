# MergeDelayErrorWith

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/MergeDelayErrorWith?tab=doc)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/merge.html)

**MergeDelayErrorWith** combines multiple Observables into one by merging their emissions.
Any error will be deferred until all observables terminate.

![MergeDelayErrorWith](../../../assets/MergeDelayErrorWith.svg?raw=true)

## Example
```go
import _ "github.com/reactivego/rx/generic"
```
Code:
```go
const _5ms = 5 * time.Millisecond
const _10ms = 10 * time.Millisecond

sourceA := CreateInt(func(N NextInt, E Error, _ Complete, _ Canceled) {
	N(1)
	E(RxError("error.sourceA"))
})

sourceB := CreateInt(func(N NextInt, _ Error, C Complete, _ Canceled) {
	time.Sleep(_5ms)
	N(0)
	time.Sleep(_10ms)
	N(2)
	C()
})

result, err := sourceA.MergeDelayErrorWith(sourceB).ToSlice()
fmt.Println(result, err)
```
Output:
```
[1 0 2] error.sourceA
```
