package rx

const RepeatCountInvalid = Error("repeat count invalid")

// Repeat creates an Observable that emits the entire source sequence multiple times.
//
// Parameters:
//   - count: Optional. The number of repetitions:
//   - If omitted: The source Observable is repeated indefinitely
//   - If 0: Returns an empty Observable
//   - If negative: Returns an Observable that emits an error
//   - If multiple count values: Returns an Observable that emits an error
//
// The resulting Observable will subscribe to the source Observable repeatedly
// each time the source completes, up to the specified count.
func Repeat[T any](count ...int) Pipe[T] {
	return func(observable Observable[T]) Observable[T] {
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
}

// Repeat emits the items emitted by the source Observable repeatedly.
//
// Parameters:
//   - count: Optional. The number of repetitions:
//   - If omitted: The source Observable is repeated indefinitely
//   - If 0: Returns an empty Observable
//   - If negative: Returns an Observable that emits an error
//   - If multiple count values: Returns an Observable that emits an error
//
// The resulting Observable will subscribe to the source Observable repeatedly
// each time the source completes, up to the specified count.
func (observable Observable[T]) Repeat(count ...int) Observable[T] {
	return Repeat[T](count...)(observable)
}
