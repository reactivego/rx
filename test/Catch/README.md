# Catch
[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/Catch#section-documentation)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/catch.html)

**Catch** recovers from an error notification by continuing the sequence without
emitting the error but by switching to the catch ObservableInt to provide items.

![Catch](../../../assets/Catch.svg?raw=true)

## Example
```go
import _ "github.com/reactivego/rx/generic"
```
Code:
```go
throw := ThrowInt(RxError("error"))
err := FromInt(1, 2, 3).ConcatWith(throw).Catch(FromInt(4, 5)).Println()
fmt.Println(err)
```
Output:
```
1
2
3
4
5
<nil>
```
