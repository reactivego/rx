# Distinct

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/Distinct?tab=doc)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/distinct.html)

**Distinct** suppress duplicate items emitted by an Observable.

The operator filters an Observable by only allowing items through that have not already been emitted.
In some implementations there are variants that allow you to adjust the criteria by which two items are
considered “distinct.” In some, there is a variant of the operator that only compares an item against its
immediate predecessor for distinctness, thereby filtering only consecutive duplicate items from the sequence.

![Distinct](../../../assets/Distinct.svg?raw=true)

## Example
```go
import _ "github.com/reactivego/rx/generic"
```
Code:
```go
FromInt(1, 2, 2, 1, 3).Distinct().Println()
```
Output:
```
1
2
3
```
