# StartWith

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/StartWith?tab=doc)
[![](../../../assets/godoc.svg?raw=true)](https://godoc.org/github.com/reactivego/rx/test/StartWith)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/startwith.html)

**StartWith** returns an observable that, at the moment of subscription, will
synchronously emit all values provided to this operator, then subscribe to
the source and mirror all of its emissions to subscribers.

![StartWith](../../../assets/StartWith.svg?raw=true)

## Example
```go
import _ "github.com/reactivego/rx/generic"
```
Code:
```go
FromInt(2, 3).StartWith(1).Println()
```
Output:
```
1
2
3
```
