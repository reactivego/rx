# PublishLast

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/PublishLast?tab=doc)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/publish.html)

**PublishLast** returns a Multicaster that shares a single subscription to the
underlying Observable containing only the last value emitted before it
completes. When the underlying Obervable terminates with an error, then
subscribed observers will receive only that error (and no value). After all
observers have unsubscribed due to an error, the Multicaster does an internal
reset just before the next observer subscribes.
