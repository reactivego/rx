package rx

import (
	"sync"
)

func WithLatestFromAll[T any](observable Observable[Observable[T]]) Observable[[]T] {
	return func(observe Observer[[]T], scheduler Scheduler, subscriber Subscriber) {
		sources := []Observable[T](nil)
		var observers struct {
			sync.Mutex
			assigned    []bool
			values      []T
			initialized int
			done        bool
		}
		makeObserver := func(sourceIndex int) Observer[T] {
			observer := func(next T, err error, done bool) {
				observers.Lock()
				defer observers.Unlock()
				if !observers.done {
					switch {
					case !done:
						if !observers.assigned[sourceIndex] {
							observers.assigned[sourceIndex] = true
							observers.initialized++
						}
						observers.values[sourceIndex] = next
						if sourceIndex == 0 && observers.initialized == len(observers.values) {
							observe(observers.values, nil, false)
						}
					case err != nil:
						observers.done = true
						var zero []T
						observe(zero, err, true)
					default:
						if sourceIndex == 0 {
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
				sources = append(sources, next)
			case err != nil:
				var zero []T
				observe(zero, err, true)
			default:
				runner := scheduler.Schedule(func() {
					if subscriber.Subscribed() {
						numSources := len(sources)
						observers.assigned = make([]bool, numSources)
						observers.values = make([]T, numSources)
						for sourceIndex, source := range sources {
							if !subscriber.Subscribed() {
								return
							}
							source.AutoUnsubscribe()(makeObserver(sourceIndex), scheduler, subscriber)
						}
					}
				})
				subscriber.OnUnsubscribe(runner.Cancel)
			}
		}
		observable(observer, scheduler, subscriber)
	}
}
