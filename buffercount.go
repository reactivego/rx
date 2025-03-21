package rx

func BufferCount[T any](observable Observable[T], bufferSize, startBufferEvery int) Observable[[]T] {
	return func(observe Observer[[]T], scheduler Scheduler, subscriber Subscriber) {
		var buffer []T
		observer := func(next T, err error, done bool) {
			switch {
			case !done:
				buffer = append(buffer, next)
				n := len(buffer)
				if n >= bufferSize {
					if n == bufferSize {
						clone := append(make([]T, 0, n), buffer...)
						observe(clone, nil, false)
					}
					if n >= startBufferEvery {
						n = copy(buffer, buffer[startBufferEvery:])
						buffer = buffer[:n]
					}
				}
			case err != nil:
				observe(nil, err, true)
			default:
				n := len(buffer)
				if 0 < n && n <= bufferSize {
					Of(buffer)(observe, scheduler, subscriber)
				} else {
					observe(nil, nil, true)
				}
			}
		}
		observable(observer, scheduler, subscriber)
	}
}
