package rx

const RepeatCountInvalid = Error("repeat count invalid")

// Repeat creates an Observable that emits a sequence of items repeatedly.
// The count is the number of times to repeat the sequence. If count is
// not provided, the sequence is repeated indefinitely.
func (observable Observable[T]) Repeat(count ...int) Observable[T] {
	if len(count) == 1 && count[0] < 0 || len(count) > 1 {
		return Throw[T](RepeatCountInvalid)
	}
	if len(count) == 1 && count[0] == 0 {
		return Empty[T]()
	}
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var repeated int
		var observer Observer[T]
		observer = func(next T, err error, done bool) {
			if !done || err != nil {
				observe(next, err, done)
			} else {
				repeated++
				if len(count) == 0 || repeated < count[0] {
					observable(observer, scheduler, subscriber)
				} else {
					var zero T
					observe(zero, nil, true)
				}
			}
		}
		observable(observer, scheduler, subscriber)
	}
}
