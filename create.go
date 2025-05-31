package rx

// Create constructs a new Observable from a Creator function.
//
// The Creator function is called repeatedly with an incrementing index value,
// and returns a tuple of (next value, error, done flag). The Observable will
// continue producing values until either:
//  1. The Creator signals completion by returning done=true
//  2. The Observer unsubscribes
//  3. The Creator returns an error (which will be emitted with done=true)
//
// This function provides a bridge between imperative code and the reactive
// Observable pattern.
func Create[T any](create Creator[T]) Observable[T] {
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		task := func(index int, again func(int)) {
			if subscriber.Subscribed() {
				next, err, done := create(index)
				if subscriber.Subscribed() {
					if !done {
						observe(next, nil, false)
						if subscriber.Subscribed() {
							again(index + 1)
						}
					} else {
						var zero T
						observe(zero, err, true)
					}
				}
			}
		}
		runner := scheduler.ScheduleLoop(0, task)
		subscriber.OnUnsubscribe(runner.Cancel)
	}
}
