# BufferTime

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/BufferTime#section-documentation)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/buffer.html)

**BufferTime** buffers the source Observable values for a specific time period and emits those as a
slice periodically in time.

Buffers values from the source during a certain period of time. It emits and resets the buffer every period.

![BufferTime](../../../assets/BufferTime.svg?raw=true)
