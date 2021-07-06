#BufferTime

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/BufferTime#section-documentation)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/buffer.html)

buffers the source Observable values for a specific time period and emits those as a
slice periodically in time.

![BufferTime](../../..assets/BufferTime.svg?raw=true)



When an observer subscribes to a `BehaviorSubject`, it begins by emitting the item most
recently emitted by the Observable part of the subject (or a seed/default
value if none has yet been emitted) and then continues to emit any other
items emitted later by the Observable part.
