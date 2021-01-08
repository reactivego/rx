# PublishBehavior

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/PublishBehavior?tab=doc)
[![](../../../assets/godoc.svg?raw=true)](https://godoc.org/github.com/reactivego/rx/test/PublishBehavior)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/publish.html)

**PublishBehavior** returns a Multicaster that shares a single subscription
to the underlying Observable returning an initial value or the last
value emitted by the underlying Observable. When the underlying
Obervable terminates with an error, then subscribed observers will
receive that error. After all observers have unsubscribed due to an error,
the Multicaster does an internal reset just before the next observer
subscribes.
