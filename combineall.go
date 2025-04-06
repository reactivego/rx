package rx

import (
	"sync"
)

func CombineAll[T any](observable Observable[Observable[T]]) Observable[[]T] {
	return func(observe Observer[[]T], scheduler Scheduler, subscriber Subscriber) {
		sources := []Observable[T](nil)
		var buffers struct {
			sync.Mutex
			assigned    []bool
			values      []T
			initialized int
			active      int
		}
		makeObserver := func(sourceIndex int) Observer[T] {
			observer := func(next T, err error, done bool) {
				buffers.Lock()
				defer buffers.Unlock()
				if buffers.active > 0 {
					switch {
					case !done:
						if !buffers.assigned[sourceIndex] {
							buffers.assigned[sourceIndex] = true
							buffers.initialized++
						}
						buffers.values[sourceIndex] = next
						if buffers.initialized == len(buffers.values) {
							observe(buffers.values, nil, false)
						}
					case err != nil:
						buffers.active = 0
						var zero []T
						observe(zero, err, true)
					default:
						if buffers.active--; buffers.active == 0 {
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
				if len(sources) == 0 {
					var zero []T
					observe(zero, nil, true)
					return
				}
				runner := scheduler.Schedule(func() {
					if subscriber.Subscribed() {
						numSources := len(sources)
						buffers.assigned = make([]bool, numSources)
						buffers.values = make([]T, numSources)
						buffers.active = numSources
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
