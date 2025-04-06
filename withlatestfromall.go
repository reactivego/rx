package rx

import (
	"sync"
)

func WithLatestFromAll[T any](observable Observable[Observable[T]]) Observable[[]T] {
	return func(observe Observer[[]T], scheduler Scheduler, subscriber Subscriber) {
		var sources []Observable[T]
		var buffers struct {
			sync.Mutex
			assigned    []bool
			values      []T
			initialized int
			done        bool
		}
		makeObserver := func(sourceIndex int) Observer[T] {
			observer := func(next T, err error, done bool) {
				buffers.Lock()
				defer buffers.Unlock()
				if !buffers.done {
					switch {
					case !done:
						if !buffers.assigned[sourceIndex] {
							buffers.assigned[sourceIndex] = true
							buffers.initialized++
						}
						buffers.values[sourceIndex] = next
						if sourceIndex == 0 && buffers.initialized == len(buffers.values) {
							observe(buffers.values, nil, false)
						}
					case err != nil:
						buffers.done = true
						var zero []T
						observe(zero, err, true)
					default:
						if sourceIndex == 0 {
							buffers.done = true
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
