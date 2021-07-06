# Publish

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/Publish#section-documentation)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/publish.html)

**Publish** returns a Multicaster for a Subject to an underlying Observable
and turns the subject into a connnectable observable.

A Subject emits to an observer only those items that are emitted by the
underlying Observable subsequent to the time of the observer subscribes.

When the underlying Obervable terminates with an error, then subscribed
observers will receive that error. After all observers have unsubscribed due
to an error, the Multicaster does an internal reset just before the next
observer subscribes. So this **Publish** operator is re-connectable, unlike
the RxJS 5 behavior that isn't. To simulate the RxJS 5 behavior use
Publish().AutoConnect(1) this will connect on the first subscription but will
never re-connect.
