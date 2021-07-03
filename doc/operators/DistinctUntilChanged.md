## DistinctUntilChanged

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/DistinctUntilChanged?tab=doc)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/distinct.html)

**DistinctUntilChanged** only emits when the current value is different from the last.

The operator only compares emitted items from the source Observable against their immediate
predecessors in order to determine whether or not they are distinct.

![DistinctUntilChanged](../../../assets/DistinctUntilChanged.svg?raw=true)

## Example
```go
import _ "github.com/reactivego/rx/generic"
```
Code:
```go
FromInt(1, 2, 2, 1, 3).DistinctUntilChanged().Println()
```
Output:
```
1
2
1
3
```
