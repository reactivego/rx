# WithLatestFromAll

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/WithLatestFromAll?tab=doc)
[![](../../../assets/rx.svg?raw=true)](https://rxjs-dev.firebaseapp.com/api/operators/withLatestFrom)

**WithLatestFromAll** flattens a higher order observable (e.g. ObservableObservable)
by subscribing to all emitted observables (ie. Observable entries) until the
source completes. It will then wait for all of the subscribed Observables to
emit before emitting the first slice. The first observable that was emitted by
the source will be used as the trigger observable. Whenever the trigger
observable emits, a new slice will be emitted containing all the latest values.

Note that any values emitted by the trigger before all other observables have emitted will
effectively be lost. The first emit will occur the first time the trigger emits after all other
observables have emitted.

![WithLatestFromAll](../../../assets/WithLatestFromAll.svg?raw=true)

## Example
```go
import _ "github.com/reactivego/rx/generic"
```
Code:
```go
a := FromInt(1, 2, 3, 4, 5).AsObservable()
b := FromString("A", "B", "C", "D", "E").AsObservable()
c := FromObservable(a, b)
c.WithLatestFromAll().Println()
```
Output:
```
[2 A]
[3 B]
[4 C]
[5 D]
```
