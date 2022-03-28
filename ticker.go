package observable

import "time"

// Ticker creates an ObservableTime that emits a sequence of timestamps after
// an initialDelay has passed. Subsequent timestamps are emitted using a
// schedule of intervals passed in. If only the initialDelay is given, Ticker
// will emit only once.
func Ticker(initialDelay time.Duration, intervals ...time.Duration) Observable[time.Time] {
	observable := func(observe Observer[time.Time], scheduler Scheduler, subscriber Subscriber) {
		i := 0
		runner := scheduler.ScheduleFutureRecursive(initialDelay, func(again func(time.Duration)) {
			if subscriber.Subscribed() {
				if i == 0 || (i > 0 && len(intervals) > 0) {
					observe(scheduler.Now(), nil, false)
				}
				if subscriber.Subscribed() {
					if len(intervals) > 0 {
						again(intervals[i%len(intervals)])
					} else {
						if i == 0 {
							again(0)
						} else {
							var zero time.Time
							observe(zero, nil, true)
						}
					}
				}
				i++
			}
		})
		subscriber.OnUnsubscribe(runner.Cancel)
	}
	return observable
}
