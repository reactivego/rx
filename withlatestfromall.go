package observable

import (
	"sync"
)

func WithLatestFromAll[T any](observable Observable[Observable[T]]) Observable[[]T] {
	return func(observe Observer[[]T], scheduler Scheduler, subscriber Subscriber) {
		observables := []Observable[T](nil)
		var observers struct {
			sync.Mutex
			assigned    []bool
			values      []T
			initialized int
			done        bool
		}
		makeObserver := func(index int) Observer[T] {
			observer := func(next T, err error, done bool) {
				observers.Lock()
				defer observers.Unlock()
				if !observers.done {
					switch {
					case !done:
						if !observers.assigned[index] {
							observers.assigned[index] = true
							observers.initialized++
						}
						observers.values[index] = next
						if index == 0 && observers.initialized == len(observers.values) {
							observe(observers.values, nil, false)
						}
					case err != nil:
						observers.done = true
						var zero []T
						observe(zero, err, true)
					default:
						if index == 0 {
							observers.done = true
							var zero []T
							observe(zero, nil, true)
						}
					}
				}
			}
			return observer
		}

		observer := func(next Observable[T], err error, done bool) {
			switch {
			case !done:
				observables = append(observables, next)
			case err != nil:
				var zero []T
				observe(zero, err, true)
			default:
				scheduler.Schedule(func() {
					if subscriber.Subscribed() {
						numObservables := len(observables)
						observers.assigned = make([]bool, numObservables)
						observers.values = make([]T, numObservables)
						for i, v := range observables {
							if !subscriber.Subscribed() {
								return
							}
							v.AutoUnsubscribe()(makeObserver(i), scheduler, subscriber)
						}
					}
				})
			}
		}
		observable(observer, scheduler, subscriber)
	}
}
