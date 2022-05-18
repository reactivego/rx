package x

import (
	"sync"
)

func CombineAll[T any](observable Observable[Observable[T]]) Observable[[]T] {
	return func(observe Observer[[]T], scheduler Scheduler, subscriber Subscriber) {
		observables := []Observable[T](nil)
		var observers struct {
			sync.Mutex
			assigned    []bool
			values      []T
			initialized int
			active      int
		}
		makeObserver := func(index int) Observer[T] {
			observer := func(next T, err error, done bool) {
				observers.Lock()
				defer observers.Unlock()
				if observers.active > 0 {
					switch {
					case !done:
						if !observers.assigned[index] {
							observers.assigned[index] = true
							observers.initialized++
						}
						observers.values[index] = next
						if observers.initialized == len(observers.values) {
							observe(observers.values, nil, false)
						}
					case err != nil:
						observers.active = 0
						var zero []T
						observe(zero, err, true)
					default:
						if observers.active--; observers.active == 0 {
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
						observers.active = numObservables
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
