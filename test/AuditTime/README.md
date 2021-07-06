# AuditTime

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/AuditTime#section-documentation)
[![](../../../assets/rx.svg?raw=true)](https://rxjs.dev/api/operators/auditTime)

**AuditTime** waits until the source emits and then starts a timer. When the timer
expires, **AuditTime** will emit the last value received from the source during the
time period when the timer was active.
