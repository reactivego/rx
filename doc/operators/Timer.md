# Timer

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/Timer?tab=doc)
[![](../../../assets/godoc.svg?raw=true)](https://godoc.org/github.com/reactivego/rx/test/Timer)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/timer.html)

**Timer** creates an Observable that emits a sequence of integers (starting at
zero) after an initialDelay has passed. Subsequent values are emitted using  a
schedule of intervals passed in. If only the initialDelay is given, **Timer** will
emit only once.
