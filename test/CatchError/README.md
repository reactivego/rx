# CatchError

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/CatchError#section-documentation)
[![](../../../assets/rx.svg?raw=true)](https://rxjs-dev.firebaseapp.com/api/operators/catchError)

**CatchError** catches errors on the Observable to be handled by returning a
new Observable or throwing an error. It is passed a selector function 
that takes as arguments err, which is the error, and caught, which is the
source observable, in case you'd like to "retry" that observable by
returning it again. Whatever observable is returned by the selector will be
used to continue the observable chain.

![CatchError](../../../assets/CatchError.svg?raw=true)

## Example
```go
import _ "github.com/reactivego/rx/generic"
```
Code:
```go
const problem = RxError("problem")

catcher := func(err error, caught ObservableInt) ObservableInt {
    if err == problem {
        return FromInt(4, 5)
    } else {
        return caught
    }
}

err := FromInt(1, 2, 3).ConcatWith(ThrowInt(problem)).CatchError(catcher).Println()

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
