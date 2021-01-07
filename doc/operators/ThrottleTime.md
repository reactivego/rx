# ThrottleTime

[![](../../../assets/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx/test/ThrottleTime?tab=doc)
[![](../../../assets/godoc.svg)](https://godoc.org/github.com/reactivego/rx/test/ThrottleTime)
[![](../../../assets/rx.svg)](https://rxjs.dev/api/operators/throttleTime)

**ThrottleTime** emits when the source emits and then starts a timer during
which all emissions from the source are ignored. After the timer expires,
**ThrottleTime** will again emit the next item the source emits, and so on.
