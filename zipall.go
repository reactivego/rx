package rx

import "sync"

func ZipAll[T any](observable Observable[Observable[T]]) Observable[[]T] {
	return func(observe Observer[[]T], scheduler Scheduler, subscriber Subscriber) {
		var sources []Observable[T]
		var buffers struct {
			sync.Mutex
			values    [][]T  // Buffer of values for each source
			completed []bool // Track which sources have completed
			active    int    // Count of active sources
		}

		makeObserver := func(sourceIndex int) Observer[T] {
			return func(next T, err error, done bool) {
				buffers.Lock()
				defer buffers.Unlock()

				if buffers.active <= 0 {
					return // Already completed or errored
				}

				switch {
				case !done:
					// Add the value to this source's buffer
					buffers.values[sourceIndex] = append(buffers.values[sourceIndex], next)

					// Check if every buffer has at least one value
					haveAllValues := true
					for _, buffer := range buffers.values {
						if len(buffer) == 0 {
							haveAllValues = false
							break
						}
					}

					// If we have values from all sources, emit the next set
					if haveAllValues {
						result := make([]T, len(buffers.values))
						for i := range buffers.values {
							result[i] = buffers.values[i][0]
							buffers.values[i] = buffers.values[i][1:] // Remove the used value
						}
						observe(result, nil, false)

						// Check if we need to complete after emission
						for i, buffer := range buffers.values {
							if len(buffer) == 0 && buffers.completed[i] {
								buffers.active = 0
								var zero []T
								observe(zero, nil, true)
								return
							}
						}
					}

				case err != nil:
					// Error terminates the entire observable
					buffers.active = 0
					buffers.values = nil
					buffers.completed = nil
					var zero []T
					observe(zero, err, true)

				default:
					// Mark this source as completed
					buffers.completed[sourceIndex] = true

					// Check if this source has an empty buffer
					if len(buffers.values[sourceIndex]) == 0 {
						// If a source completes with no buffered values, we can't emit more values
						buffers.active = 0
						buffers.values = nil
						buffers.completed = nil
						var zero []T
						observe(zero, nil, true)
						return
					}

					// Otherwise decrement active count and check if all sources are done
					buffers.active--
					if buffers.active == 0 {
						buffers.values = nil
						buffers.completed = nil
						var zero []T
						observe(zero, nil, true)
					}
				}
			}
		}

		observer := func(next Observable[T], err error, done bool) {
			switch {
			case !done:
				sources = append(sources, next)
			case err != nil:
				var zero []T
				observe(zero, err, true)
			default:
				scheduler.Schedule(func() {
					if subscriber.Subscribed() {
						numSources := len(sources)
						if numSources == 0 {
							var zero []T
							observe(zero, nil, true)
							return
						}

						buffers.values = make([][]T, numSources)
						buffers.completed = make([]bool, numSources)
						buffers.active = numSources

						for i, source := range sources {
							if !subscriber.Subscribed() {
								return
							}
							source.AutoUnsubscribe()(makeObserver(i), scheduler, subscriber)
						}
					}
				})
			}
		}
		observable(observer, scheduler, subscriber)
	}
}
