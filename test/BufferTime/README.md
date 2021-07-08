# BufferTime

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/BufferTime#section-documentation)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/buffer.html)

**BufferTime** buffers the source Observable values for a specific time period and emits those as a
slice periodically in time.

Buffers values from the source during a certain period of time. It emits and resets the buffer every period.

![BufferTime](../../../assets/BufferTime.svg?raw=true)

## Example
```go
import _ "github.com/reactivego/rx/generic"
```
Code:
```go
const ms = time.Millisecond
source := TimerInt(0*ms, 100*ms).Take(4).ConcatMap(func(i int) Observable {
	switch i {
	case 0:
		return From("a", "b")
	case 1:
		return From("c", "d", "e")
	case 3:
		return From("f", "g")
	}
	return Empty()
})
source.BufferTime(100 * ms).Println()
```
Output:
```
[a b]
[c d e]
[]
[f g]
```
